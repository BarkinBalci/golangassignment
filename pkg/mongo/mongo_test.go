package mongo

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const testDBName = "test_db"

func setupTestClient(t *testing.T) *Client {
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("DB_NAME", testDBName)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create MongoDB client: %v", err)
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := client.client.Database(testDBName).Drop(ctx)
		if err != nil {
			t.Errorf("Failed to drop test database: %v", err)
		}

		err = client.client.Disconnect(ctx)
		if err != nil {
			t.Errorf("Failed to disconnect from MongoDB: %v", err)
		}
	})

	return client
}

func insertTestDocuments(t *testing.T, client *Client, records []interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.db.Collection("records")

	_, err := collection.InsertMany(ctx, records)
	if err != nil {
		t.Fatalf("Failed to insert test documents: %v", err)
	}
}

func TestFetchData(t *testing.T) {
	client := setupTestClient(t)

	testRecords := []interface{}{
		bson.M{"key": "A", "createdAt": time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "counts": []int{1, 2, 3}},
		bson.M{"key": "B", "createdAt": time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC), "counts": []int{4, 5}},
		bson.M{"key": "C", "createdAt": time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC), "counts": []int{1}},
	}
	insertTestDocuments(t, client, testRecords)

	tests := []struct {
		name      string
		startDate string
		endDate   string
		minCount  int
		maxCount  int
		expected  []Record
		wantErr   bool
	}{
		{
			name:      "ValidDateRangeAndCount",
			startDate: "2024-01-15",
			endDate:   "2024-01-25",
			minCount:  5,
			maxCount:  10,
			expected: []Record{
				{Key: "B", CreatedAt: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC), TotalCount: 9},
			},
			wantErr: false,
		},
		{
			name:      "NoMatchingRecords",
			startDate: "2024-01-15",
			endDate:   "2024-01-25",
			minCount:  10,
			maxCount:  20,
			expected:  []Record{},
			wantErr:   false,
		},
		{
			name:      "InvalidStartDate",
			startDate: "invalid-date",
			endDate:   "2024-01-25",
			minCount:  0,
			maxCount:  100,
			expected:  nil,
			wantErr:   true,
		},
		{
			name:      "InvalidEndDate",
			startDate: "2024-01-15",
			endDate:   "invalid-date",
			minCount:  0,
			maxCount:  100,
			expected:  nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := client.FetchData(tt.startDate, tt.endDate, tt.minCount, tt.maxCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.expected) {
					t.Errorf("FetchData() got %d records, want %d records", len(got), len(tt.expected))
					return
				}
				for i := range got {
					if got[i].Key != tt.expected[i].Key || got[i].TotalCount != tt.expected[i].TotalCount || !got[i].CreatedAt.Equal(tt.expected[i].CreatedAt) {
						t.Errorf("FetchData() got [%d] = %+v, want %+v", i, got[i], tt.expected[i])
					}
				}
			}
		})
	}
}
