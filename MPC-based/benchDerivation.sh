#!/bin/bash

# Step 1: Compile the program
./compile.py sha512

# Step 2: Generate secure channels for all parties
Scripts/setup-ssl.sh 128

# Step 3: Define models and parties to run
models=("shamir" "mal-shamir")
parties=(3 10 16 32 64 128)
thresholds=(1 4 7 15 31 63)

# Function to run a model with specified number of parties
run_model() {
    local model=$1
    local num_parties=$2
	local thres=$3
    local output_file="output-${model}.out"

    # Run the model with the given number of parties
    case $model in
        "shamir")
            PLAYERS=$num_parties Scripts/shamir.sh sha512 -T $thres >> $output_file
            ;;
        "mal-shamir")
            PLAYERS=$num_parties Scripts/mal-shamir.sh sha512 -T $thres >> $output_file
            ;;
        # "atlas")
        #     Scripts/atlas.sh sha512 -N $num_parties -T $thres >> $output_file
        #     ;;
        *)
            echo "Unknown model: $model"
            ;;
    esac
}

# Loop through each model and each party count, then run and save the output
for model in "${models[@]}"; do
    for (( i=0; i<${#parties[@]}; i++ )); do
        party_count=${parties[$i]}
        threshold=${thresholds[$i]}
        echo "Running ${model} with ${party_count} parties and threshold ${threshold}..."
        run_model $model $party_count $threshold
    done
done
