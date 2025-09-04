package utils

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Orange-Health/citadel/common/constants"
)

func ConvertStringToUint(str string) uint {
	convertedValue, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return uint(convertedValue)
}

func ConvertStringToUintSlice(str string) []uint {
	strArray := strings.Split(str, ",")
	var uintArray []uint
	for _, str := range strArray {
		convertedValue, _ := strconv.ParseUint(str, 10, 64)
		uintArray = append(uintArray, uint(convertedValue))
	}
	return uintArray
}

func ConvertStringToStringSlice(str string) []string {
	str = strings.TrimSpace(str)
	if str == "" {
		return []string{}
	}
	strArray := strings.Split(str, ",")
	for i, s := range strArray {
		strArray[i] = strings.TrimSpace(s)
	}
	return strArray
}

func ConvertUintToString(element uint) string {
	return strconv.FormatUint(uint64(element), 10)
}

func ConvertUintSliceToString(uintArray []uint) string {
	var strArray []string
	for _, value := range uintArray {
		strArray = append(strArray, strconv.Itoa(int(value)))
	}
	return strings.Join(strArray, ",")
}

func ConvertStringSliceToString(slice []string) string {
	return strings.Join(slice, ",")
}

func ConvertStringToInt(str string) int {
	convertedValue, _ := strconv.ParseInt(str, 10, 64)
	return int(convertedValue)
}

func ConvertStringToFloat32ForAbnormality(s string) float32 {
	s = strings.Trim(s, "<=> ")
	s = strings.ReplaceAll(s, " ", "")
	return ConvertStringToFloat32(s)
}

func ConvertStringToFloat32(s string) float32 {
	f, _ := strconv.ParseFloat(s, 32)
	return float32(f)
}

func ConvertIntSliceToUintSlice(intSlice []int) []uint {
	uintSlice := make([]uint, len(intSlice))
	for index, value := range intSlice {
		uintSlice[index] = uint(value)
	}
	return uintSlice
}

func GetNonEmptyString(s1, s2 string) string {
	if s1 != "" {
		return strings.TrimSpace(s1)
	}
	return strings.TrimSpace(s2)
}

func GetDobByYearsMonthsDays(years, months, days int) time.Time {
	currentTime := time.Now()
	dayStartingTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())

	return dayStartingTime.AddDate(-years, -months, -days)
}

func GetAgeYearsMonthsAndDaysFromDob(dob time.Time) (uint, uint, uint) {
	currentTime := time.Now()
	dayStartingTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())

	ageYears := dayStartingTime.Year() - dob.Year()
	ageMonths := int(dayStartingTime.Month()) - int(dob.Month())
	ageDays := dayStartingTime.Day() - dob.Day()

	if ageDays < 0 {
		previousMonth := dayStartingTime.AddDate(0, -1, 0)
		daysInPrevMonth := time.Date(previousMonth.Year(), previousMonth.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		ageDays += daysInPrevMonth
		ageMonths--
	}

	if ageMonths < 0 {
		ageYears--
		ageMonths += 12
	}

	return uint(ageYears), uint(ageMonths), uint(ageDays)
}

func GetAgeYearsFromDob(dob time.Time) uint {
	currentTime := time.Now()
	dayStartingTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())

	ageYears := dayStartingTime.Year() - dob.Year()
	if dayStartingTime.Month() < dob.Month() || (dayStartingTime.Month() == dob.Month() && dayStartingTime.Day() < dob.Day()) {
		ageYears--
	}

	return uint(ageYears)
}

func GetStringRequestIdWithoutStringPart(omsRequestId string) string {
	if omsRequestId == "" {
		return ""
	}

	// Using regex to find the numeric part of the OMS order ID
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(omsRequestId)
	return match
}

func GetUintRequestIdWithoutStringPart(omsRequestId string) uint {
	omsRequestIdStr := GetStringRequestIdWithoutStringPart(omsRequestId)
	omsRequestIdUint, _ := strconv.ParseUint(omsRequestIdStr, 10, 64)
	return uint(omsRequestIdUint)
}

func GetStringOrderIdWithoutStringPart(omsOrderId string) string {
	if omsOrderId == "" {
		return ""
	}

	// Using regex to find the numeric part of the OMS order ID
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(omsOrderId)
	return match
}

func GetUintOrderIdWithoutStringPart(omsOrderId string) uint {
	omsOrderIdStr := GetStringOrderIdWithoutStringPart(omsOrderId)
	omsOrderIdUint, _ := strconv.ParseUint(omsOrderIdStr, 10, 64)
	return uint(omsOrderIdUint)
}

func GetUintTestIdWithoutStringPart(omsTestId string) uint {
	omsTestIdStr := GetStringOrderIdWithoutStringPart(omsTestId)
	omsTestIdUint, _ := strconv.ParseUint(omsTestIdStr, 10, 64)
	return uint(omsTestIdUint)
}

func GetCurrentFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	fullName := function.Name()

	// Split the full name by "/" to get the function name with package path
	parts := strings.Split(fullName, "/")

	// The actual function name with package will be the last part
	fullName = parts[len(parts)-1]

	// Split the function name with package path by "." to get the function name
	parts = strings.Split(fullName, ".")
	functionName := parts[len(parts)-1]

	return functionName
}

func CalculatePatientDob(years, months, days int) string {
	if years == 0 && months == 0 && days == 0 {
		return ""
	}

	loc, _ := time.LoadLocation(constants.LocalTimeZoneLocation)
	dateToday := time.Now().In(loc)
	dateToday = dateToday.AddDate(-int(years), -int(months), -int(days))
	dateOfBirth := dateToday.Format(constants.DateLayout)
	return dateOfBirth
}

func GetGenderConstant(s string) string {
	switch strings.ToLower(s) {
	case "male":
		return "M"
	case "female":
		return "F"
	}

	return "O"
}

func GetGenderConstantForPatient(genderUint uint) string {
	switch genderUint {
	case 1:
		return "male"
	case 2:
		return "female"
	}

	return "other"
}

func GetFileExtension(fileName string) string {
	parts := strings.Split(fileName, ".")
	return parts[len(parts)-1]
}

func MaskPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) < 6 {
		return phoneNumber
	}

	return fmt.Sprintf("%sXXXX%s", phoneNumber[:3], phoneNumber[len(phoneNumber)-3:])
}

func ToTitleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(strings.ToUpper(s)[0]) + strings.ToLower(s)[1:]
}

func GetLocalLocation() (*time.Location, error) {
	return time.LoadLocation(constants.LocalTimeZoneLocation)
}

func GetSalutationByGender(gender string) string {
	if strings.ToLower(gender) == "female" {
		return "Ms."
	}

	return "Mr."
}

func GetPatientIdByOrderIdAndPatientId(omsOrderId string, patientId string) string {
	if patientId != "" {
		return patientId
	}

	return fmt.Sprintf("OH%s", omsOrderId)
}

func GetDeduplicationIdForTestEvent(testId string) string {
	return fmt.Sprint(testId, ":", GetCurrentTimeInMilliseconds())
}

func GetOmsBaseDomain(cityCode string) string {
	cityCode = strings.ToLower(cityCode)
	baseDomain := constants.OmsBaseDomain
	if cityCode != "" {
		baseDomain = strings.ReplaceAll(baseDomain, "city_code", cityCode)
	}
	return baseDomain
}

func getSuperlabInfoScreenUrl(orderId string) string {
	return fmt.Sprintf("%s/#/info?orderId=%s", constants.SuperlabBaseUrl, orderId)
}
