@echo on
echo %cd%

echo "构建 %cd%\%1%.dll"

go build -gcflags "all=-N -l" -ldflags "-s -w" -o %1_tmp.dll -buildmode=c-archive %1.go
upx %1_tmp.dll -o %1.dll
pause