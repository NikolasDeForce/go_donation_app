# Golang Donation App 'DonateON'

# Приложение для сбора донатов от зрителей на трансляциях. Реализовано 3 микросервиса, которые собраны в один монолит.
# Первый сервис - основной сайт DonateON, на котором показывают преимущества приложения и предлагают зарегистрироваться. Регистрационные данные сохраняются в БД PostgreSQL. При регистрации генерируется персональный токен для обращения по API и присваивается клиенту. Клиент может узнать свой токен по запросу /api/{login}/{password}/mytoken
# Второй сервис реализует прием донатов стримеру от зрителя. В полях зритель указывает нужные данные о себе, стримере, вводит сумму доната, сообщение и платежные данные. Ссылка: http://localhost:8010/donation
# Третий сервис реализует обработку запросов от стримеров через API. Список донатов от зрителей /api/{token}/donates 

Старт - `docker-compose up` потом - `go run main.go`

Если проблема с PostgreSQL, то нужно переместить create_db.sql на машину командой - `psql -U postgres postgres -h 127.0.0.1 < create_db.sql`
Либо в Docker руками перекинуть в папку и проинициализировать командой - `psql -U postgres postgres < create_db.sql`

Главная страница http://localhost:8010/
![Alt text](prew/prew1.png?raw=true "Main")

Страница отправки доната стримеру http://localhost:8010/donation
![Alt text](prew/prew2.png?raw=true "Donate")

Страница API запроса для получения токена /api/{login}/{password}/mytoken
![Alt text](prew/prew3.png?raw=true "Token")

Страница API запроса для получения списка донатов по токену стримера /api/{token}/donates
![Alt text](prew/prew4.png?raw=true "List")
