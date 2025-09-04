package calculationsService

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"

	commonConstants "github.com/Orange-Health/citadel/common/constants"
	commonStructs "github.com/Orange-Health/citadel/common/structures"
	commonUtils "github.com/Orange-Health/citadel/common/utils"
	commonModels "github.com/Orange-Health/citadel/models"
)

type CalculationsServiceInterface interface {
	GetModifiedValue(ctx context.Context, formula string, decimals uint, dependentMasterInvestigationIds []uint,
		masterInvestigationIdToInvestigationResultMap map[uint]commonModels.InvestigationResult) (
		string, *commonStructs.CommonError)
}

func (calculationsService *CalculationsService) GetModifiedValue(ctx context.Context, formula string, decimals uint,
	dependentMasterInvestigationIds []uint,
	masterInvestigationIdToInvestigationResultMap map[uint]commonModels.InvestigationResult) (
	string, *commonStructs.CommonError) {

	if formula == "" {
		return "", &commonStructs.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    commonConstants.ERROR_NO_FORMULA_FOUND,
		}
	}

	for _, masterInvestigationId := range dependentMasterInvestigationIds {
		investigationResult, ok := masterInvestigationIdToInvestigationResultMap[masterInvestigationId]
		if !ok {
			return "", &commonStructs.CommonError{
				StatusCode: http.StatusInternalServerError,
				Message:    commonConstants.ERROR_NO_INVESTIGATION_RESULTS_FOUND,
			}
		}

		formula = strings.ReplaceAll(formula, fmt.Sprintf("{{%d}}", masterInvestigationId),
			investigationResult.InvestigationValue)
	}

	expression, err := govaluate.NewEvaluableExpression(formula)
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
			"formula": formula,
		}, err)
		return "", &commonStructs.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    commonConstants.ERROR_IN_CALCULATING_RESULT,
		}
	}

	result, err := expression.Evaluate(map[string]interface{}{})
	if err != nil {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
			"formula": formula,
		}, err)
		return "", &commonStructs.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    commonConstants.ERROR_IN_CALCULATING_RESULT,
		}
	}

	resultString := ""
	resultFloat, ok := result.(float64)
	if ok {
		resultString = roundFloat(resultFloat, decimals)
	} else {
		commonUtils.AddLog(ctx, commonConstants.ERROR_LEVEL, commonUtils.GetCurrentFunctionName(), map[string]interface{}{
			"formula": formula,
		}, errors.New(commonConstants.ERROR_INVALID_RESULT))
		return "", &commonStructs.CommonError{
			StatusCode: http.StatusInternalServerError,
			Message:    commonConstants.ERROR_INVALID_RESULT,
		}
	}

	return resultString, nil
}

func roundFloat(val float64, precision uint) string {
	ratio := math.Pow(10, float64(precision))
	result := math.Round(val*ratio) / ratio
	return strconv.FormatFloat(result, 'f', -1, 64)
}
