// Package unit provides the Unit type, representing a value with no information.
//
// Unit is useful in generic contexts where a function signature requires a type parameter
// but the value carries no meaningful data (similar to void in languages like Java or Haskell's ()).
//
// Example usage with Result:
//
//	func deleteUser(id string) result.Result[unit.Unit] {
//	    if err := db.Delete(id); err != nil {
//	        return result.Err[unit.Unit](err)
//	    }
//	    return result.Ok(unit.Unit{})
//	}
//
// Or with Writer for logging:
//
//	effect.NewWriter(result.Ok(unit.Unit{}), []string{"user deleted"})
package unit
