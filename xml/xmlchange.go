package xml

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Patient struct {
	Cardnumber int    `json:"cardnumber"`
	Surname    string `json:"surname"`
	Name       string `json:"Name"`
	Secondname string `json:"secondname"`
}

func openTemplate(name string) ([]byte, error) {
	filename := name
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// хранить данные о пациентах лучше в sql таблице, но для данного тестового задания я ограничился только хэш-таблицей
// чтобы реализовать занесение в шаблон разных пациентов

func InsertInXml(template string, id int) string {
	data, err := openTemplate(template)
	if err != nil {
		log.Fatal(err)
	}
	PatientMap := map[int]Patient{30: {123456, "Ivanov", "Ivan", "Ivanovich"}, 31: {111222, "Sidorov", "Sidor", "Sidorovich"}, 32: {228282, "John", "Dow", "Chingizidovich"}}
	p := PatientMap[id]
	// находим блок в который нужно вставить значение по тегу и выделяем его в отдельный срез байт
	// реализация проста, по тегу ищем начало начало блока и определяем индекс его первго и последнего элемента
	// в этом блоке заменяем "_"
	// CARDNUMBER
	firstCardNumIndex := strings.Index(string(data), "CARDNUMBER")
	endCardNumIndex := strings.Index(string(data[firstCardNumIndex:]), "</ns1:text>")
	cardNumSlice := data[firstCardNumIndex : firstCardNumIndex+endCardNumIndex]
	// в срезе заменяем нижный пробел на нужное значение
	cardNumString := strings.Replace(string(cardNumSlice), "_", strconv.Itoa(p.Cardnumber), 1)
	// объединяем все в новый срез
	newData := string(data[:firstCardNumIndex]) + cardNumString + string(data[firstCardNumIndex+endCardNumIndex:])
	// SURNAME
	firstSurNameIndex := strings.Index(newData, "surname")
	endSurNameIndex := strings.Index(newData[firstSurNameIndex:], "</ns1:text>")
	surnameSlice := newData[firstSurNameIndex : firstSurNameIndex+endSurNameIndex]
	surnameString := strings.Replace(surnameSlice, "_", " "+p.Surname, 1)
	newData = newData[:firstSurNameIndex] + surnameString + newData[firstSurNameIndex+endSurNameIndex:]
	// NAME
	firstNameIndex := strings.Index(newData, "=\"name")
	endNameIndex := strings.Index(newData[firstNameIndex:], "</ns1:text>")
	nameSlice := newData[firstNameIndex : firstNameIndex+endNameIndex]
	nameString := strings.Replace(nameSlice, "_", " "+p.Name, 1)
	newData = newData[:firstNameIndex] + nameString + newData[firstNameIndex+endNameIndex:]
	// SECONDNAME
	firstSNameIndex := strings.Index(newData, "secondname")
	endSNameIndex := strings.Index(newData[firstSNameIndex:], "</ns1:text>")
	secondNameSlice := newData[firstSNameIndex : firstSNameIndex+endSNameIndex]
	secondNameString := strings.Replace(secondNameSlice, "_", " "+p.Secondname, 1)
	newData = newData[:firstSNameIndex] + secondNameString + newData[firstSNameIndex+endSNameIndex:]
	// Записываем измененную в 4 частях строку в файл c именем в формате YYYY MM DD hh mm ss
	currentTime := time.Now()
	time := currentTime.Format("2006-01-02 15:04:05")
	time = strings.Replace(time, ":", ".", -1)
	filename := time + ".doc"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(newData)
	return filename
}
