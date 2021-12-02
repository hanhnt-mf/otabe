package controller

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "otabe/pb"
	"strconv"
)

const (
	maxPageLimit = 200
)

func ValidateLocationConditions(location *pb.SearchLocationConditions) *pb.ValidationErrorDetails {
	if location == nil {
		return nil
	}
	if location.Long == nil || location.Lat == nil || location.Distance == nil {
		return &pb.ValidationErrorDetails{
			Fields: []*pb.ValidationTargetField{
				{
					Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_LONGITUDE,
					FieldLocation: []*pb.InvalidFieldLocation{
						{Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_CONDITIONS},
						{Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_SEARCH_LOCATIONS_CONDITIONS},
					},
				},
				{
					Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_LATITUDE,
					FieldLocation: []*pb.InvalidFieldLocation{
						{Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_CONDITIONS},
						{Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_SEARCH_LOCATIONS_CONDITIONS},
					},
				},
				{
					Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_DISTANCE,
					FieldLocation: []*pb.InvalidFieldLocation{
						{Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_CONDITIONS},
						{Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_SEARCH_LOCATIONS_CONDITIONS},
					},
				},
			},
			Reason: pb.InvalidReasonType_INVALID_REASON_TYPE_MUST_BE_SET_TOGETHER,
			Description: "longitude, latitude and distance fields must be set together",
		}
	}

	return nil
}

func ValidatePaging(paging *pb.Paging) []*pb.ValidationErrorDetails {
	if paging ==  nil {
		return nil
	}
	var validationErrorDetailsList []*pb.ValidationErrorDetails

	if paging.PageLimit > maxPageLimit {
		validationErrorDetailsList = append(validationErrorDetailsList, &pb.ValidationErrorDetails{
			Fields: []*pb.ValidationTargetField{
				{
					Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_MAX_CONTENTS_PER_PAGE,
					FieldLocation: []*pb.InvalidFieldLocation{
						{Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_PAGING},
					},
				},
			},
			Reason: pb.InvalidReasonType_INVALID_REASON_TYPE_MUST_BE_LESS_THAN_OR_EQUAL_TO,
			ReasonOptions: map[string]string{
				"max_value": strconv.Itoa(maxPageLimit),
			},
			Description: fmt.Sprintf("max contents per page must be less than or equal to %d", maxPageLimit),
		})
	}

	if paging.PageLimit < 1 {
		validationErrorDetailsList = append(validationErrorDetailsList, &pb.ValidationErrorDetails{
			Fields: []*pb.ValidationTargetField{
				{
					Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_MAX_CONTENTS_PER_PAGE,
					FieldLocation: []*pb.InvalidFieldLocation{
						{Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_PAGING},
					},
				},
			},
			Reason: pb.InvalidReasonType_INVALID_REASON_TYPE_MUST_BE_GREATER_THAN_OR_EQUAL_TO,
			ReasonOptions: map[string]string{
				"min_value": "1",
			},
			Description: fmt.Sprintf("max contents per page must be greater than or equal to 1"),
		})
	}

	if paging.PageNumber < 1 {
		validationErrorDetailsList = append(validationErrorDetailsList, &pb.ValidationErrorDetails{
			Fields: []*pb.ValidationTargetField{
				{
					Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_PAGE_NUMBER,
					FieldLocation: []*pb.InvalidFieldLocation{
						{Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_PAGING},
					},
				},
			},
			Reason: pb.InvalidReasonType_INVALID_REASON_TYPE_MUST_BE_GREATER_THAN_OR_EQUAL_TO,
			ReasonOptions: map[string]string{
				"min_value": "1",
			},
			Description: "page number must be greater than or equal to 1",
		})
	}

	return validationErrorDetailsList
}

func RestaurantNotFound() error {
	var validationErrorDetailsList []*pb.ValidationErrorDetails

	validationErrorDetails := &pb.ValidationErrorDetails{
		Fields: []*pb.ValidationTargetField{
			{
				Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_RESTAURANT_ID,
				FieldLocation: []*pb.InvalidFieldLocation{
					{Field: pb.InvalidFieldType_INVALID_FIELD_TYPE_RESTAURANT_PARAMS},
				},
			},
		},
		Reason: pb.InvalidReasonType_INVALID_REASON_TYPE_RESTAURANT_NOT_FOUND,
		ReasonOptions: map[string]string{
			"restaurant_id": "1",
		},
		Description: "restaurant with id is not found",
	}
	validationErrorDetailsList = append(validationErrorDetailsList, validationErrorDetails)

	return NewInvalidArgumentErrorWithDetails(validationErrorDetailsList)
}

func ValidateListRestaurantsRequest(req *pb.ListRestaurantsRequest) error {
	var validationErrorDetailsList []*pb.ValidationErrorDetails

	validationErrorDetails := ValidateLocationConditions(req.Location)
	if validationErrorDetails != nil {
		validationErrorDetailsList = append(validationErrorDetailsList, validationErrorDetails)
	}
	// need to handle sort by column

	validationErrorDetailsList = append(validationErrorDetailsList, ValidatePaging(req.Paging)...)
	if len(validationErrorDetailsList) == 0 {
		return nil
	}
	return NewInvalidArgumentErrorWithDetails(validationErrorDetailsList)
}



func NewInvalidArgumentErrorWithDetails(validationErrorDetailsList []*pb.ValidationErrorDetails) error {
	if validationErrorDetailsList == nil {
		return nil
	}
	st := status.New(codes.InvalidArgument, "validation error")
	var err error
	for _, details := range validationErrorDetailsList {
		st, err = st.WithDetails(details)
		if err != nil {
			// 想定外のエラーなので、黙って client に err を返すより panic して rollbar に出力したほうが良い
			panic(fmt.Sprintf("Unexpected error when appending details to error: %v", err))
		}
	}

	return st.Err()
}