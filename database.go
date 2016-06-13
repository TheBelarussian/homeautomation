package main

// Work in progress exporting db calls from model to this file

// PUBLIC functions for database management
/*
func getBucketDataAsBytes(bucketName string) []byte {

	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		b.ForEach(func(k, v []byte) error {
			d := Device{}
			// Decode json value of cell to Device struct type
			err := json.Unmarshal(v, &d)
			if err != nil {
				return nil
			}
			devices = append(devices, d)
			fmt.Println(k, d)

			return nil
		})

		return nil
	})
}*/
