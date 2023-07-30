My solution to [Fetch's receipt processor challenge](https://github.com/fetch-rewards/receipt-processor-challenge). Directions assume Go is already installed per project description. It is also assumed that Go is correctly installed on the path.

### Directions (Windows):
1: Clone repository to directory of your choosing\
2: Open a terminal and path to project directory\
3: Run "go get ." to download dependancies\
4: Run "go run ." to start web service\
5: Web service is now started and ready to receive CRUD requests at the specified endpoints. Program defaults to http://localhost:8080. If you wish to change this default, change the "hosturl" string in view.go.
