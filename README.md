# Simple API notes

Простое приложение. Хотел отработать работу с маршрутизатором Chi, хотел построить согласно [архитектуре](https://habr.com/ru/articles/269589/) entity-action-adaptor-service. И пощупать автогенерацию swagger документации из комментариев в коде сервиса http.

## References

- [router CHI](https://go-chi.io/#/README)
- [github chi](https://github.com/go-chi/chi/tree/master/_examples)
- 
- [driver postgresql](https://pkg.go.dev/github.com/lib/pq#section-readme)
- [github driver postgresql](https://github.com/lib/pq/tree/master)
- [simple example for use this driver](https://metanit.com/go/tutorial/10.3.php)
- [how to use array for postgresql](https://www.opsdash.com/blog/postgres-arrays-golang.html)

- [docker image postgres](https://hub.docker.com/_/postgres)

- [autogenerate swagger from comments](https://github.com/swaggo/swag/tree/master?tab=readme-ov-file#descriptions-over-multiple-lines)
  - [example](https://github.com/swaggo/swag/blob/master/example/celler/controller/examples.go)
  - [another example](https://github.com/swaggo/swag/blob/master/example/celler/main.go)
  - [and another](https://github.com/swaggo/http-swagger/blob/master/example/go-chi/main.go)

- [cors chi middleware](https://pkg.go.dev/github.com/go-chi/cors)
- [github cors package](https://github.com/go-chi/cors)

- [GitHub REST API Documentation](https://docs.github.com/en/rest?apiVersion=2022-11-28#cross-origin-resource-sharing)