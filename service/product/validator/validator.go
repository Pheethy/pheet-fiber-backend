package validate

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pheet-fiber-backend/models"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Validation struct{}

var validate = validator.New()
func ValidateStruct(pro *models.Products) []*models.Element {
    var errors = make([]*models.Element,0)
    err := validate.Struct(pro)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            var element = new(models.Element)
            element.FailedField = err.StructNamespace()
            element.Tag = err.Tag()
            element.Value = err.Param()
            errors = append(errors, element)
        }
    }
    return errors
}

func (v Validation) ValidateCreate() fiber.Handler {
	var next fiber.Handler
	return func(c *fiber.Ctx) error {
		var body = c.Body()
		var a interface{}
		json.Unmarshal(body, &a)

		var params = a.(map[string]interface{})
		/* key params */
		key := "name"
		name, nameOK := params[key]
		log.Println("name:", name)
		if !nameOK {
			return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("%s: was missing on body", key))
			// return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s: was missing on body", key).Error())
		}
		if name != nil && name != "" {
			if err := validation.Validate(name, validation.By(models.ValidateTypeString)); err != nil {
				return c.Status(http.StatusBadRequest).SendString(fmt.Sprintf("%s: %s", key, err.Error()))
			}
		}
		// key = "alias_name"
		// aliasName, aliasNameOK := params[key]
		// if !aliasNameOK {
		// 	return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s: was missing on body", key).Error())
		// }
		// if aliasName != nil && aliasName != "" {
		// 	if err := validation.Validate(aliasName, validation.By(helper.ValidateTypeString)); err != nil {
		// 		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s: %s", key, err.Error()).Error())
		// 	}
		// }
		// //validate AgencyNo
		// key = "agency_no"
		// agencyNo, agencyNoOK := params[key]
		// if !agencyNoOK {
		// 	return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s: was missing on body", key).Error())
		// }
		// if agencyNo != nil && agencyNo != "" {
		// 	if err := validation.Validate(agencyNo, validation.By(helper.ValidateTypeString)); err != nil {
		// 		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s: %s", key, err.Error()).Error())
		// 	}
		// }
		// //validate AgencyType
		// key = "agency_type"
		// agencyType, agencyTypeOK := params[key]
		// if !agencyTypeOK {
		// 	return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s: was missing on body", key).Error())
		// }
		// if agencyType != nil && agencyType != "" {
		// 	if err := validation.Validate(agencyType, validation.By(helper.ValidateTypeString)); err != nil {
		// 		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s: %s", key, err.Error()).Error())
		// 	}
		// }
		// key = "private_tel_no"
		// privateTelNo, privateTelNoOK := params[key]
		// if !privateTelNoOK {
		// 	return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s: was missing on body", key).Error())
		// }
		// if privateTelNo != nil && privateTelNo != "" {
		// 	if err := validation.Validate(privateTelNo, validation.By(helper.ValidateTypeString)); err != nil {
		// 		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s: %s", key, err.Error()).Error())
		// 	}
		// }
		// key = "parent_ids"
		// parentIds, parentIdsOK := params[key]
		// if !parentIdsOK {
		// 	return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s: was missing on body", key).Error())
		// }
		// if parentIds != nil && reflect.TypeOf(parentIds).Kind() == reflect.Slice {
		// 	if len(parentIds.([]interface{})) > 0 {
		// 		if err := validation.Validate(parentIds, validation.Each(is.UUIDv4)); err != nil {
		// 			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("%s: %s", key, err.Error()).Error())
		// 		}
		// 	}
		// }

		return next(c)
	}
}