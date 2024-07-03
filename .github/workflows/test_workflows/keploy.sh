
chmod +x ${GITHUB_WORKSPACE}/scripts/keploy-initialize-env.sh

# Initialize the environment
sh ${GITHUB_WORKSPACE}/scripts/keploy-initialize-env.sh

# # Get into virtual environment
# source .venv/bin/activate
# echo "Activated virtual environment"

export KEPLOY_API_KEY=TL0Z2j0AxW649wE4Tg==

# Get the Keploy binary
curl --silent -o keployE --location https://keploy-enterprise.s3.us-west-2.amazonaws.com/releases/0.7.11/enterprise_linux_amd64
sudo chmod a+x keployE && sudo mkdir -p /usr/local/bin && sudo mv keployE /usr/local/bin


chmod +x ${GITHUB_WORKSPACE}/scripts/keploy_local_server.sh
sudo -E env PATH="$PATH" /usr/local/bin/keployE test -c "${GITHUB_WORKSPACE}/scripts/keploy_local_server.sh" --delay 20 --apiTimeout 300 --freezeTime --generateGithubActions=false
echo "Keploy started in test mode"

all_passed=true

# Loop through test sets
for i in {0..0}
do
    # Define the report file for each test set
    report_file="./keploy/reports/test-run-0/test-set-$i-report.yaml"

    # Extract the test status
    test_status=$(grep 'status:' "$report_file" | head -n 1 | awk '{print $2}')

    # Print the status for debugging
    echo "Test status for test-set-$i: $test_status"

    # Check if any test set did not pass
    if [ "$test_status" != "PASSED" ]; then
        all_passed=false
        echo "Test-set-$i did not pass."
        break # Exit the loop early as all tests need to pass
    fi
done

# Check the overall test status and exit accordingly
if [ "$all_passed" = true ]; then
    echo "All tests passed"
    exit 0
else
    exit 1
fi