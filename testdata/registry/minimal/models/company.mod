name: Company
fields:
  ID:
    type: AutoIncrement
  Name:
    type: String
  TaxID:
    type: String
identifiers:
  primary: ID
  name: Name
related:
  Person:
    type: HasMany
  MainContact:
    type: ForOne
    aliased: Contact
  MailingContact:
    type: ForOne
    aliased: Contact
  Note:
    type: HasManyPoly
    through: Commentable
    aliased: Comment
