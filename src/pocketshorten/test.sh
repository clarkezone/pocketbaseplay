	podman manifest exists localhost/pocketshortenp:latest
	if [ $? -eq 0 ];
	    	then
		echo "Manifest exists"
		podman manifest rm localhost/pocketshortenp:latest
		false
	else
		echo "Manifest doesn't exist"
	fi
