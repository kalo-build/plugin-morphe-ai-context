name: Person
fields:
  ID:
    type: AutoIncrement
  FirstName:
    type: String
  LastName:
    type: String
  Nationality:
    type: Nationality
identifiers:
  primary: ID
  name:
    - FirstName
    - LastName
related:
  ContactInfo:
    type: HasOne
  Company:
    type: ForOne
  WorkContact:
    type: ForOne
    aliased: Contact
  PersonalContact:
    type: ForOne
    aliased: Contact
  Note:
    type: HasManyPoly
    through: Commentable
    aliased: Comment
