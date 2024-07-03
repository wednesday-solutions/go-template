echo
echo "Copying .env.keploy.linux to .env.local"
cp .env.keploy.linux .env.local

echo "Enabling Pre Commit Hooks"
pre-commit install
echo

opentelemetry-bootstrap --action=install