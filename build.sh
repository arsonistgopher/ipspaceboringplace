Data1='{"cuidsid":'
Data2='}'
Data3=$Data1\"$1\"$Data2

echo "---- Build Script for Basic Automation Demo ---"

RESULT=$(curl -s -o /dev/null -w "%{http_code}" -X POST -H 'Content-Type: application/json' -d $Data3  localhost:1323/vars)
echo "Created resources on AutoPAM with ID "$1
echo "AutoPAM result is "$RESULT

if [ $RESULT -eq 201 ]; then
	git clone https://github.com/arsonistgopher/L3VPNAssets.git
	mv L3VPNAssets $1
	echo "Moved directory to "$1
	cp ~/.ssh/id_auto ./$1
	echo "Copied automation SSH key to "$1
	mkdir -p ./$1/host_vars/$2
	cp ./$1/utils/credentials.yaml ./$1/host_vars/$2/credentials.yaml
	echo "Created host_var directory and inserted credentials"
   else
	echo "Error: Cannot get variables"
fi
	
