1. Clone the repository to your local machine.
2. Apply the yamls under manifests folder.
3. Exec into the mysql pod and run: 
```
# mysql -u root --password=<redacted>
# GRANT ALL PRIVILEGES ON <database_name>.* TO '<user>'@'<rest-api_pod_IP>' IDENTIFIED BY '<password>';
# FLUSH PRIVILEGES;
# Run the sql scripts 
```
4. kubectl port-forward -n restapi svc/restapi 8080
5. Run GET/POST request from Postman. 
