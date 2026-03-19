package command

import (
	"bufio"
	"errors"
	"fmt"
	"labs/pkg/models"
	"labs/pkg/worker"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Command struct {
	CommType string `json:"type"`
	Args     string `json:"text"`
}

func ParseCommand(line string) (Command, error) {
	if line == "" {
		return Command{}, errors.New("Пустая команда")
	}
	parts := strings.SplitN(line, " ", 2)
	if len(parts) < 2 {
		return Command{}, errors.New("Неполная команда")
	}
	cmdType := strings.ToUpper(parts[0])
	args := strings.TrimSpace(parts[1])

	switch cmdType {
	case "ADD":
		return Command{CommType: cmdType, Args: args}, nil
	case "DEL":
		return Command{CommType: cmdType, Args: args}, nil
	case "SAVE":
		return Command{CommType: cmdType, Args: args}, nil
	default:
		return Command{}, errors.New("Неизвестная комманда")
	}
}

func ExecuteSave(args string, products *[]models.Product) error {
	filename := strings.TrimSpace(args)
	if filename == "" {
		return errors.New("не указано имя файла для сохранения")
	}
	return worker.WriteFile(filename, worker.CreateProductString(*products))
}

func ExecuteAdd(args string, products *[]models.Product, dataFilePath string) error {
	var dateStr string
	var p models.Product

	_, err := fmt.Sscanf(args, "%q %s %d %s", &p.Name, &p.Category, &p.Count, &dateStr)
	if err != nil {
		return errors.New("Некорректная строка")
	}

	p.Date, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		return errors.New("Ошибка парсинга даты")
	}
	p.ID = uuid.New().String()

	*products = append(*products, p)

	return worker.WriteFile(dataFilePath, worker.CreateProductString(*products))
}

var conditionRegex = regexp.MustCompile(`^(\w+)\s*(<=|>=|==|!=|<|>)\s*(.+)$`)

func ExecuteDel(args string, products *[]models.Product, dataFilePath string) error {
	field, op, valueStr, err := parseCondition(args)
	if err != nil {
		return err
	}

	predicate, err := buildPredicate(field, op, valueStr)
	if err != nil {
		return err
	}

	*products = filter(*products, predicate)
	return worker.WriteFile(dataFilePath, worker.CreateProductString(*products))
}

func parseCondition(args string) (field, op, value string, err error) {
	matches := conditionRegex.FindStringSubmatch(args)
	if len(matches) != 4 {
		return "", "", "", errors.New("неверный формат условия DEL")
	}
	return matches[1], matches[2], strings.TrimSpace(matches[3]), nil
}

func buildPredicate(field, op, valueStr string) (func(models.Product) bool, error) {
	switch field {
	case "name", "category":
		return buildStringPredicate(field, op, valueStr)
	case "count":
		return buildIntPredicate(op, valueStr)
	case "date":
		return buildDatePredicate(op, valueStr)
	default:
		return nil, errors.New("неизвестное поле условия")
	}
}

func buildStringPredicate(field, op, valueStr string) (func(models.Product) bool, error) {
	if !strings.HasPrefix(valueStr, `"`) || !strings.HasSuffix(valueStr, `"`) {
		return nil, errors.New("значение для поля name должно быть в двойных кавычках")
	}
	strVal := strings.Trim(valueStr, `"`)

	if op != "==" && op != "!=" {
		return nil, errors.New("для строкового поля %s допустимы только операции == и !=")
	}

	return func(p models.Product) bool {
		var fieldVal string
		if field == "name" {
			fieldVal = p.Name
		} else {
			fieldVal = p.Category
		}
		if op == "==" {
			return fieldVal == strVal
		}
		return fieldVal != strVal
	}, nil
}

func buildIntPredicate(op, valueStr string) (func(models.Product) bool, error) {
	val, err := strconv.Atoi(valueStr)
	if err != nil {
		return nil, errors.New("значение count должно быть целым числом")
	}

	ops := map[string]func(int, int) bool{
		"<":  func(a, b int) bool { return a < b },
		">":  func(a, b int) bool { return a > b },
		"<=": func(a, b int) bool { return a <= b },
		">=": func(a, b int) bool { return a >= b },
		"==": func(a, b int) bool { return a == b },
		"!=": func(a, b int) bool { return a != b },
	}
	compare, ok := ops[op]
	if !ok {
		return nil, errors.New("неподдерживаемый оператор для count")
	}

	return func(p models.Product) bool {
		return compare(p.Count, val)
	}, nil
}

func buildDatePredicate(op, valueStr string) (func(models.Product) bool, error) {
	dateVal, err := time.Parse("2006-01-02", valueStr)
	if err != nil {
		return nil, errors.New("неверный формат даты, ожидается YYYY-MM-DD")
	}

	ops := map[string]func(time.Time, time.Time) bool{
		"<":  func(a, b time.Time) bool { return a.Before(b) },
		">":  func(a, b time.Time) bool { return a.After(b) },
		"<=": func(a, b time.Time) bool { return a.Before(b) || a.Equal(b) },
		">=": func(a, b time.Time) bool { return a.After(b) || a.Equal(b) },
		"==": func(a, b time.Time) bool { return a.Equal(b) },
		"!=": func(a, b time.Time) bool { return !a.Equal(b) },
	}
	compare, ok := ops[op]
	if !ok {
		return nil, errors.New("неподдерживаемый оператор для date")
	}

	return func(p models.Product) bool {
		return compare(p.Date, dateVal)
	}, nil
}

func filter(products []models.Product, pred func(models.Product) bool) []models.Product {
	result := make([]models.Product, 0, len(products))
	for _, p := range products {
		if !pred(p) {
			result = append(result, p)
		}
	}
	return result
}

func LoadCommandsFromFile(path string) []Command {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Command{}
		}
		log.Printf("Ошибка открытия файла команд %s: %v", path, err)
		return []Command{}
	}
	defer file.Close()

	var records []Command
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			log.Printf("Некорректная строка команды (пропущена): %s", line)
			continue
		}
		cmdType := parts[0]
		text := parts[1]

		record := Command{
			Args:     text,
			CommType: cmdType,
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения файла команд: %v", err)
	}
	return records
}

func SaveCommandsToFile(filePath string, list []Command) error {

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, rec := range list {
		line := fmt.Sprintf("%s %s\n", rec.CommType, rec.Args)
		if _, err := file.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}
