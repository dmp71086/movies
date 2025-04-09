# Урок 12

## Tracing

1. http://localhost:16686/

## Prometheus

1. http://localhost:9090/

## Grafana

1. http://localhost:3000/
1. See more: https://grafana.com/grafana/dashboards/

## Pprof

1. `brew install graphviz`
1. http://localhost:8000/debug/pprof/


### Profile (CPU)
1. `go tool pprof http://localhost:8000/debug/pprof/profile?seconds=30s`
1. `go tool pprof -http=localhost:4000 <path to profile>`

### Heap
1. `go tool pprof http://localhost:8000/debug/pprof/heap`
1. `go tool pprof -http=localhost:4000 <path to profile>`

### Mutex
1. `go tool pprof http://localhost:8000/debug/pprof/mutex?debug=1`


## Полезные материалы

### Форматирование и стилизация кода

- https://talks.golang.org/2014/names.slide
- https://golang.org/doc/effective_go.html#names
- https://blog.golang.org/package-names
- https://github.com/golang/go/wiki/CodeReviewComments
- [Руководство по стилизации кода для пакетов Golang](https://rakyll.org/style-packages/) (rakyll/JBD)

### Стандартный макет Go проекта

- https://github.com/golang-standards/project-layout/blob/master/README_ru.md#%D0%B4%D0%B8%D1%80%D0%B5%D0%BA%D1%82%D0%BE%D1%80%D0%B8%D0%B8-go

### Наименование и организация пакетов и структуру кода

- [GopherCon EU 2018: Peter Bourgon - Best Practices for Industrial Programming](https://www.youtube.com/watch?v=PTE4VJIdHPg)
- [GopherCon Russia 2018: Ashley McNamara + Brian Ketelsen - Go best practices.](https://www.youtube.com/watch?v=MzTcsI6tn-0)
- [GopherCon 2017: Edward Muller - Go Anti-Patterns](https://www.youtube.com/watch?v=ltqV6pDKZD8)
- [GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)


### Чистая архитектура

- https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- https://amitshekhar.me/blog/go-backend-clean-architecture
- https://habr.com/ru/companies/oleg-bunin/articles/516500/
- https://github.com/evrone/go-clean-template/tree/master
- https://github.com/bxcodec/go-clean-arch


### Go Style Guides

- https://go.dev/doc/effective_go
- https://github.com/uber-go/guide/blob/master/style.md#uber-go-style-guide
- https://google.github.io/styleguide/go/
- https://google.github.io/styleguide/go/best-practices
