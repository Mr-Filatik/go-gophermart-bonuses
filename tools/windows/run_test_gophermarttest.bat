echo off

cd ..\..\
cd cmd\accrual\
echo Run build for accrual application
go build -o accrual.exe main.go
echo Build for accrual application is done

cd ..\..\
cd cmd\gophermart\
echo Run build for gophermart application
go build -o gophermart.exe main.go
echo Build for gophermart application is done

cd ..\..\

echo Run gophermarttest
gophermarttest-windows-amd64 ^
  -test.v -test.run=^TestGophermart$ ^
  -gophermart-binary-path=cmd/gophermart/gophermart ^
  -gophermart-host=localhost ^
  -gophermart-port=8080 ^
  -gophermart-database-uri="postgres://user:password@localhost:5432/postgres?sslmode=disable" ^
  -accrual-binary-path=cmd/accrual/accrual ^
  -accrual-host=localhost ^
  -accrual-port=8081 ^
  -accrual-database-uri="postgres://user:password@localhost:5432/postgres?sslmode=disable"

pause

echo on