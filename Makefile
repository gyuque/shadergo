GOC     = 6g
GOL     = 6l
PRODUCT = main

$(PRODUCT): main.6
	$(GOL) -o $(PRODUCT) $^

main.6: main.go canvas.6 file.6 shader.6 gyu3d.6
	$(GOC) main.go

canvas.6: canvas.go shader.6 gyu3d.6
	$(GOC) canvas.go

file.6: file.go
	$(GOC) file.go

gyu3d.6: gyu3d.go
	$(GOC) gyu3d.go
	
shader.6: shader.go gyu3d.6
	$(GOC) shader.go

.PHONY: clean
clean:
	$(RM) $(PRODUCT) main.6 canvas.6 file.6 gyu3d.6