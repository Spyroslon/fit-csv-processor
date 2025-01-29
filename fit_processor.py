import fitparse
from datetime import datetime
import pandas as pd
import sys
import os

def load_fit_file(fit_file, selected_records=None):
    """Load a FIT file and extract relevant records."""
    try:
        if not os.path.isfile(fit_file):
            raise FileNotFoundError(f"FIT file {fit_file} not found")

        with open(fit_file, 'rb') as f:
            fit_data = fitparse.FitFile(f)
            records = []

            for record in fit_data.get_messages():
                if selected_records is None or record.name in selected_records:
                    record_dict = {}

                    for data_field in record:
                        value = data_field.value
                        
                        if isinstance(value, datetime):
                            value = value.isoformat()  # Convert datetime to ISO 8601 string

                        record_dict[data_field.name] = value

                    records.append(record_dict)

        return pd.DataFrame(records)

    except Exception as e:
        raise RuntimeError(f"Error loading FIT file: {e}")

def process_fit_file(fit_file, output_file, selected_records=None, rounded_timestamp_seconds=0):
    """Process a FIT file into a summarized CSV."""
    try:
        df = load_fit_file(fit_file, selected_records)
        
        # Convert timestamp to datetime for proper grouping
        df['timestamp'] = pd.to_datetime(df['timestamp'])

        df_cleaned = df.groupby('timestamp').mean().reset_index().dropna()

        if rounded_timestamp_seconds > 0:
            try:
                df_cleaned = df_cleaned.copy()
                df_cleaned['rounded_timestamp'] = df_cleaned['timestamp'].dt.round(f'{rounded_timestamp_seconds}s') 
                
                # Group by rounded timestamp and calculate the mean
                final_df = df_cleaned.groupby('rounded_timestamp').mean().reset_index()

                # Drop the 'timestamp' column
                if 'timestamp' in final_df.columns:
                    final_df = final_df.drop('timestamp', axis=1)

            except Exception as e:
                raise RuntimeError(f"Error processing timestamps: {e}")
        else:
            final_df = df_cleaned
        
        if not final_df.empty:
            final_df.to_csv(output_file, encoding='utf-8', index=False)
            return True
        else:
            print(f"No valid data found in {fit_file}.")
            return False

    except Exception as e:
        print(f"Error converting FIT file: {e}")
            
def main():

    selected_records = ["record", "timestamp", "distance", "enhanced_altitude", "enhanced_speed", "gps_accuracy", "position_lat", "position_long", "speed", "heart_rate"]
    rounded_timestamp_seconds = 0

    args = sys.argv[1:]
    if not args:
        print("Usage: pass fit files in arguments")
        return
        
    for filename in args:
        if not filename.endswith('.fit'):
            print(f"Skipping non-FIT file: {filename}")
            continue
            
        try:
            output_file = f"{os.path.splitext(filename)[0]}_summary.csv"
            success = process_fit_file(filename, output_file, selected_records, rounded_timestamp_seconds)
            if success:
                print(f"Successfully processed: {filename}")
            else:
                print(f"Error processing: {filename}")
                
        except Exception as e:
            print(f"Error processing file {filename}: {str(e)}")

if __name__ == "__main__":
    main()
