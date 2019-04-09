package jwt_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/dgrijalva/jwt-go/v4/test"
)

const (
	expireAt = 100
	nowInt   = 200
)

var claimsTestData = []struct {
	name   string
	need   string
	claims jwt.Claims
}{
	{
		name:   "Map claims expired json.Number",
		need:   "Given the need to test validating a MapClaims past its json.Number expire-at time.",
		claims: jwt.MapClaims{"foo": "bar", "exp": json.Number(fmt.Sprintf("%v", time.Unix(expireAt, 0).Unix()))},
	},
	{
		name:   "Map claims expired float64",
		need:   "Given the need to test validating a MapClaims past its float64 expire-at time.",
		claims: jwt.MapClaims{"foo": "bar", "exp": float64(time.Unix(expireAt, 0).Unix())},
	},
	{
		name:   "StandardClaims expired",
		need:   "Given the need to test validating a StandardClaims past its expire-at time.",
		claims: jwt.StandardClaims{ExpiresAt: jwt.NewTime(expireAt)},
	},
}

func TestClaimValidExpired(t *testing.T) {
	for _, data := range claimsTestData {
		t.Log(data.name)
		t.Logf("\t%s", data.need)
		name := data.name
		t.Logf("\t\tValidate the Claims with exp as %v at time %v", nowInt, expireAt)
		test.At(time.Unix(nowInt, 0), func() {
			err := data.claims.Valid()
			t.Log("\t\t\tExpect an error that includes the expired by 1m40s information")
			if err == nil {
				t.Errorf("[%v] Expecting error.  Didn't get one.", name)
			} else {
				ve := err.(*jwt.ValidationError)
				// compare the bitfield part of the error
				if e := ve.Errors; e != jwt.ValidationErrorExpired {
					t.Errorf("[%v] Errors don't match expectation.  %v != %v", name, e, jwt.ValidationErrorExpired)
				}
				switch vi := ve.Inner.(type) {
				default:
					expectedErrorStr := "token is expired by 1m40s"
					if fmt.Sprint(ve.Inner.Error()) != expectedErrorStr {
						t.Errorf("[%v] Errors inner text is not as expected. \"%v\" is not \"%v\"", name, ve.Inner, expectedErrorStr)
					}
				case *jwt.ExpiredError:
					if vi.ExpiredBy != 100*time.Second {
						t.Errorf("[%v] ExpiredError.ExpiredBy %v is not %v\n", name, vi.ExpiredBy, 100*time.Second)
					}
					if !reflect.DeepEqual(vi.Claims, data.claims) {
						t.Errorf("[%v] Claims did not get set in expired error \"%v\"\n", name, vi.Error())
					}
					if vi.Error() != "Token is expired" {
						t.Errorf("[%v] Error message is not as expected \"%v\"\n", name, vi.Error())
					}

				}
			}
		})
	}
}
