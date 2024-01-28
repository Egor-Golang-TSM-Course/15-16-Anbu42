# Пример использования API с использованием curl

Примеры curl-команд для взаимодействия с HTTP-сервером, реализующим API для работы с задачами и пользователями.

## Задачи (Tasks)

### CreateTask
```bash
curl -X POST -H "Content-Type: application/json" -d '{"id": 1, "name": "New Task"}' http://localhost:8080/tasks
```

### GetTasks:
```bash
curl http://localhost:8080/tasks
```

### GetTask:
```bash
curl http://localhost:8080/tasks/1
```

### UpdateTask:
```bash
curl -X PUT -H "Content-Type: application/json" -d '{"id": 1, "name": "Updated Task"}' http://localhost:8080/tasks/1
```

### DeleteTask:
```bash
curl -X DELETE http://localhost:8080/tasks/1
```

## Пользователи (Users)

### CreateUser:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"id": 1, "name": "New User"}' http://localhost:8080/users
```

### GetUsers:
```bash
curl http://localhost:8080/users
```

### GetUser:
```bash
curl http://localhost:8080/users/1
```

### UpdateUser:
```bash
curl -X PUT -H "Content-Type: application/json" -d '{"id": 1, "name": "Updated User"}' http://localhost:8080/users/1
```

### DeleteUser:
```bash
curl -X DELETE http://localhost:8080/users/1
```