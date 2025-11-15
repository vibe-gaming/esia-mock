package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// UserData содержит моковые данные пользователя
type UserData struct {
	OID         string
	FirstName   string
	LastName    string
	MiddleName  string
	BirthDate   string
	Gender      string
	SNILS       string
	INN         string
	Email       string
	Mobile      string
	Trusted     bool
	Verified    bool
	Citizenship string
	Status      string
}

// Cache хранит связь между номерами телефонов и моковыми данными пользователей
type Cache struct {
	users map[string]*UserData
	mu    sync.RWMutex
}

// New создает новый кеш
func New() *Cache {
	return &Cache{
		users: make(map[string]*UserData),
	}
}

// GetOrCreate возвращает существующие данные для телефона или создает новые
func (c *Cache) GetOrCreate(phoneNumber string) *UserData {
	c.mu.RLock()
	user, exists := c.users[phoneNumber]
	c.mu.RUnlock()

	if exists {
		return user
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Проверяем еще раз (double-check locking)
	if user, exists := c.users[phoneNumber]; exists {
		return user
	}

	// Создаем уникальные данные на основе номера телефона
	user = c.generateUserData(phoneNumber)
	c.users[phoneNumber] = user

	return user
}

// generateUserData генерирует уникальные моковые данные на основе номера телефона
func (c *Cache) generateUserData(phoneNumber string) *UserData {
	// Используем хеш телефона для генерации детерминированных, но уникальных данных
	hash := sha256.Sum256([]byte(phoneNumber))
	hashStr := hex.EncodeToString(hash[:])

	// Генерируем OID из первых 10 символов хеша
	oidNum := parseHexToNumber(hashStr[:10])
	oid := fmt.Sprintf("%d", 1000000000+oidNum%1000000000)

	// Имена из разных наборов в зависимости от хеша
	firstNames := []string{"Иван", "Петр", "Сергей", "Александр", "Дмитрий", "Андрей", "Михаил", "Алексей", "Николай", "Владимир"}
	lastNames := []string{"Иванов", "Петров", "Сидоров", "Смирнов", "Кузнецов", "Попов", "Васильев", "Павлов", "Соколов", "Михайлов"}
	middleNames := []string{"Иванович", "Петрович", "Сергеевич", "Александрович", "Дмитриевич", "Андреевич", "Михайлович", "Алексеевич", "Николаевич", "Владимирович"}

	firstNameIdx := int(hash[0]) % len(firstNames)
	lastNameIdx := int(hash[1]) % len(lastNames)
	middleNameIdx := int(hash[2]) % len(middleNames)

	// Генерируем дату рождения (годы 1970-2000)
	year := 1970 + (int(hash[3]) % 31)
	month := 1 + (int(hash[4]) % 12)
	day := 1 + (int(hash[5]) % 28) // безопасное значение для любого месяца
	birthDate := fmt.Sprintf("%02d.%02d.%d", day, month, year)

	// Пол (M или F)
	gender := "M"
	if hash[6]%2 == 0 {
		gender = "F"
		// Для женщин используем женские окончания отчеств
		femaleMiddleNames := []string{"Ивановна", "Петровна", "Сергеевна", "Александровна", "Дмитриевна", "Андреевна", "Михайловна", "Алексеевна", "Николаевна", "Владимировна"}
		middleNames = femaleMiddleNames
		middleNameIdx = int(hash[2]) % len(middleNames)

		// Женские фамилии с окончанием -ова/-ева
		lastNames = []string{"Иванова", "Петрова", "Сидорова", "Смирнова", "Кузнецова", "Попова", "Васильева", "Павлова", "Соколова", "Михайлова"}
		lastNameIdx = int(hash[1]) % len(lastNames)

		// Женские имена
		firstNames = []string{"Мария", "Анна", "Елена", "Ольга", "Татьяна", "Наталья", "Ирина", "Светлана", "Екатерина", "Юлия"}
		firstNameIdx = int(hash[0]) % len(firstNames)
	}

	// Генерируем SNILS (11 цифр)
	snilsNum := parseHexToNumber(hashStr[10:21])
	snils := fmt.Sprintf("%011d", snilsNum%100000000000)

	// Генерируем ИНН (12 цифр)
	innNum := parseHexToNumber(hashStr[21:33])
	inn := fmt.Sprintf("%012d", innNum%1000000000000)

	// Email на основе имени и хеша (транслитерация + lowercase)
	emailHash := hashStr[33:40]
	email := strings.ToLower(fmt.Sprintf("%s.%s.%s@example.com",
		transliterate(firstNames[firstNameIdx]),
		transliterate(lastNames[lastNameIdx]),
		emailHash))

	return &UserData{
		OID:         oid,
		FirstName:   firstNames[firstNameIdx],
		LastName:    lastNames[lastNameIdx],
		MiddleName:  middleNames[middleNameIdx],
		BirthDate:   birthDate,
		Gender:      gender,
		SNILS:       snils,
		INN:         inn,
		Email:       email,
		Mobile:      phoneNumber,
		Trusted:     hash[7]%2 == 0, // 50% trusted
		Verified:    hash[8]%3 != 0, // ~66% verified
		Citizenship: "RUS",
		Status:      "REGISTERED",
	}
}

// parseHexToNumber парсит hex строку в uint64
func parseHexToNumber(hexStr string) uint64 {
	num, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0
	}
	return num
}

// transliterate простая транслитерация русских имен в латиницу
func transliterate(name string) string {
	translit := map[rune]string{
		// Заглавные
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D",
		'Е': "E", 'Ё': "Yo", 'Ж': "Zh", 'З': "Z", 'И': "I",
		'Й': "Y", 'К': "K", 'Л': "L", 'М': "M", 'Н': "N",
		'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T",
		'У': "U", 'Ф': "F", 'Х': "H", 'Ц': "Ts", 'Ч': "Ch",
		'Ш': "Sh", 'Щ': "Sch", 'Ъ': "", 'Ы': "Y", 'Ь': "",
		'Э': "E", 'Ю': "Yu", 'Я': "Ya",
		// Строчные
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d",
		'е': "e", 'ё': "yo", 'ж': "zh", 'з': "z", 'и': "i",
		'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n",
		'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t",
		'у': "u", 'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch",
		'ш': "sh", 'щ': "sch", 'ъ': "", 'ы': "y", 'ь': "",
		'э': "e", 'ю': "yu", 'я': "ya",
	}

	result := ""
	for _, r := range name {
		if t, ok := translit[r]; ok {
			result += t
		} else {
			result += string(r)
		}
	}
	return result
}

// GetAll возвращает все сохраненные данные (для отладки)
func (c *Cache) GetAll() map[string]*UserData {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]*UserData, len(c.users))
	for k, v := range c.users {
		result[k] = v
	}
	return result
}

// Count возвращает количество сохраненных пользователей
func (c *Cache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.users)
}
