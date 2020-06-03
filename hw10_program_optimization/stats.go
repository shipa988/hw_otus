package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

const (
	sabaken = "@"
)

var (
	user      User
	userCount int
)

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

//GetDomainStat read json data from r and returns  DomainStat of domain and first error
func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %s", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

//getUsers read json data from r and returns []User array and error
func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	i := 0
	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), &user)
		if err != nil {
			return
		}
		result[i] = user
		i++
	}

	if err = scanner.Err(); err != nil {
		return
	}
	userCount = i // for truncate users in countDomains
	return
}

//countDomains read []User array and returns DomainStat of domain and first error
func countDomains(u users, domain string) (result DomainStat, err error) {
	result = make(DomainStat)
	lendomain := len(domain)
	var email string
	var emailParts []string
	for _, user = range u[:userCount] {
		if strings.LastIndex(user.Email, domain) == len(user.Email)-lendomain { //replace regex
			email = strings.ToLower(user.Email)
			emailParts = strings.Split(email, sabaken)
			if err == nil && len(emailParts) <= 1 || len(emailParts) > 2 {
				err = fmt.Errorf("invalid email %v", user.Email)
				continue
			}
			result[emailParts[1]]++
		}
	}
	return result, err
}
