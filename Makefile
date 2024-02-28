air:
	air 

signin:
	curl -X POST http://localhost:8000/signin -d '{"username":"user2","password":"password2"}' -v

validate:
	curl -X GET http://localhost:8000/welcome -v --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIyIiwiZXhwIjoxNzA5MTA0OTY5fQ.m9PvBKgAAhfFx0CpdTSV_Lap4FGWk072ZXis2NNFrpI"

refresh:
	curl -X GET http://localhost:8000/refresh -v --cookie "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIyIiwiZXhwIjoxNzA5MTA0OTY5fQ.m9PvBKgAAhfFx0CpdTSV_Lap4FGWk072ZXis2NNFrpI"