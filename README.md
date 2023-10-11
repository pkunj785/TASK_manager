Welcome to Task manager 

- project is about to ustilize task summery via simple crud operation. 
- prerequisition is to use Golang for Rest http service to handel user task's data.

# Used tech or tools for project as :

- Golang Chi router becauser it's minimal router lib to can max use of    StdLib of golang.
- As Database use mongoDB, reason behind this it's all project feature under document paradigm and not 
  need much any relation based operation and mongoDB provide adequate read/write operation. 
- for middleware verifications use StdLib of Golang.
- Use own KV local memory base cache service for handle User Identity for specific auth purpose.
- makeFile for project build

# Project Operation :

- first install depdency after fork using command 'make dep'
- for run server use run command 'make run'
- for cleaning binaries and project dep run command 'make clean'

# api endspoints: 

- via 'localhost:8080/registor' for userid and password registeration with POST method.
- '/login' for login for access userdata with POST method.
- after that operation for tasks as :
  1) create task -> "/api/createTask" with POST method
  1) get all tasks -> "/api/getTasks" with GET method
  1) task completion -> "/api/task/{id}" with PUT method
  1) task undo -> "/api/undoTask/{id}" with PUT method
  1) delete one task -> "/api/deleteTask/{id}" with DELETE method
  1) delete all tasks -> "/api/deleteAllTasks" with DELETE method


If optionally need to deploy server for "railway.app" could helpful.

# notes:

- All important api local variable and credential is inside .env file