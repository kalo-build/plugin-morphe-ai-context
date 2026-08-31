name: ContactInfo
fields:
  ID:
    type: AutoIncrement
  Email:
    type: String
identifiers:
  primary: ID
  email: Email
related:
  Person:
    type: ForOne
  RelatedContact:
    type: ForOne
    aliased: Contact
