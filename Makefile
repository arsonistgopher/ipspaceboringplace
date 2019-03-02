build:
	$(info ************  Building ************)
	./DBGenSrv/CreateDB/CreateDB
	cp automationvars.db ./DBGenSRv/CreateDB/automationvars.db
	mv automationvars.db ./DBGenSRv/ServeDB/automationvars.db

clean:
	$(info ************  Cleaning ************)
	rm ./DBGenSrv/CreateDB/automationvars.db 
	rm ./DBGenSrv/ServeDB/automationvars.db
	rm -rf L3VPNAssets

create:
	$(info ************  Creating ************)
	sh build.sh "$(id)" "${device}"
	./BuildAnsibleVars/BuildAnsibleVars "$(id)" ./"$(id)"/host_vars/"$(device)"/"$(device)".yaml


deploy:
	$(info ************  Deploying ************)
	cd ./"$(id)"/; \ansible-playbook build_pe.pb --extra-var target="$(device)";
	cd ./"$(id)"/; \ansible-playbook documents.pb --extra-var target="$(device)" --extra-var service="$(id)";


