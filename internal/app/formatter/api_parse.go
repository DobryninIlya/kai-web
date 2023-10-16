package formatter

import (
	"errors"
	"fmt"
	"main/internal/app/model"
	"strings"
)

const (
	sentenceMinimumCapacity = 15
	blockMaximumLength      = 200
)

var (
	ErrHeaderNotFound       = errors.New("header not found")
	ErrBodyContainsBadWords = errors.New("body contains bad words")
	ErrDescriptionNotFound  = errors.New("description not found")
)

var (
	badWords = []string{"цена", "стоимость", "оплата", "рублей", "рубля",
		"убить", "убивать", "убил", "смерть", "избил",
		"секс", "порно", "блять", "сука", "тварь", "мразь", "уебо",
		"путин", "украина", "казино", "приват", "казино"}
)

// GetHeader получат заголовок для поста - первое предложение
func GetHeader(body string) (string, error) {
	blocks := strings.Split(body, "\n")
	sentences := make([][]string, len(blocks), len(blocks))
	for i := 0; i < len(sentences); i++ {
		sentences[i] = append(sentences[i], strings.Split(blocks[i], ". ")...)
	}
	var header string
	for i, block := range sentences {
		for _, sentence := range block {
			var breakFlag bool
			for i, ch := range sentence {
				if i > 5 { // Если символ не в начале, то выходим
					break
				}
				if ch == '#' || ch == '@' || ch == '*' { // Если есть такие символы, скорее всего хэштэг, игнорируем
					breakFlag = true
					break
				}
			}
			if breakFlag {
				continue
			}
			if len(sentence) > sentenceMinimumCapacity { // Игнорируем предложения, меньше фиксированной длины
				header = sentence
				return header, nil
			}
		}
		if i > 3 {
			return "", ErrHeaderNotFound
		}
	}
	return "", ErrHeaderNotFound
}

// ValidateText проверяет тело текста на наличие запрещенных слов по словарю
func ValidateText(body string) error {
	body = strings.ToLower(body)
	for _, word := range badWords {
		if strings.Contains(body, word) {
			return ErrBodyContainsBadWords
		}
	}
	return nil
}

// GetImages формирует html структуру с картинками, идущими подряд
func GetImages(attachments []model.Attachment) string {
	var result string
	template := "<img class=\"content-img\" src=\"%v\">\n"
	for i, attachment := range attachments {
		if i == 0 { // Игнорируем, т.к первая пошла в превью
			continue
		}
		if attachment.Type != "photo" { // Игнорируем не картинки
			continue
		}
		result += fmt.Sprintf(template, attachment.Photo.Sizes[len(attachment.Photo.Sizes)-1].Url)
	}
	return result
}

// GetTagsInText возвращает тэг из текста
func GetTagsInText(body string) string {
	var result string
	var started bool
	for _, ch := range body {
		if ch == '#' {
			started = true
			result += string(ch)
			continue
		}
		if !(ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= 'А' && ch <= 'Я' || ch >= 'а' && ch <= 'я' || ch >= '0' && ch <= '9' || ch == '@') {
			return result
		}
		if started {
			result += string(ch)
		}

	}
	return result
}

// truncateString обрезает строку до length
func truncateString(s string, length int) string {
	runes := []rune(s)
	if len(runes) > length {
		return string(runes[:length])
	}
	return s
}

// GetDescription ищет текст, подходящий в качестве текстового превью
func GetDescription(body string, header string) (string, error) {
	blocks := strings.Split(body, "\n")
	sentences := make([][]string, len(blocks), len(blocks))
	for i := 0; i < len(sentences); i++ {
		sentences[i] = append(sentences[i], strings.Split(blocks[i], ". ")...)
	}
	for _, block := range blocks {
		var breakFlag bool
		if len(block) < sentenceMinimumCapacity {
			continue
		}
		for i, ch := range block {
			if i > 5 { // Если символ не в начале, то выходим
				break
			}
			if ch == '#' || ch == '@' || ch == '*' { // Если есть такие символы, скорее всего хэштэг, игнорируем
				breakFlag = true
				break
			}
		}
		if breakFlag { // Если вначале хэштэг, то игнорируем
			continue
		}
		if block == header { // Если подходящее описание совпадает с заголовком, то игнорируем
			continue
		}
		if len(block) < blockMaximumLength {
			return block, nil
		} else {
			return truncateString(block, blockMaximumLength), nil
		}

	}
	return "", ErrDescriptionNotFound
}
