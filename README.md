# Go Lang Rest API

After the two required APIs which it consumes are executed, it should be expected that the API runs correctly.

In order to call the APIs, you have to register the host name of ingress for each API. In order to get correct response, you have to register ip address of Ingress controller of Country Api in etc/resolve.conf on the pod. Additionally, I could not run two  APIs on my local computer because of lack of the resources. After I send to the endpoint of airplane API, all ingress controller can not respond any call. That's why, I read data of airplane on a json file.

//you can run this command in the following line to get in the pod.

kubectl exec -it countryairportlist-api-deployment-59c7f896cd-6wg9h -n countryairportlist-api -- /bin/sh

//you need to update the files below.
- vi /etc/resolv.conf 

- vi /etc/hosts 

Finally, you can find the source code which I implemented to consume the country and airport APIs.


Example:

You can call use these links in the following line to get response message.

- http://countryairportlist-api.info/countryairportsummary?runwayminimum=2

- http://countryairportlist-api.info/countryairportsummary?runwayminimum

- http://countryairportlist-api.info/countryairportsummary/2

- http://countryairportlist-api.info/countryairportsummary
