# cosys - content builder
This module is responsible for handling the creation of new content types.

### Documentation:

```yaml
paths:
  /admin/get/{contentTypeName}:
    get:
      summary: Returns a content type schema by name
      parameters:
        - in: path
          name: contentTypeName
          description: Name of the content type schema to return
          required: true
          schema:
            type: string
  /admin/build/{contentTypeName}:
    post:
      summary: Create a new content type
      requestBody:
        content:
          application/x-yaml:
            schema:
              $ref: '#/components/schemas/Schema'
        required: true
components:
  schemas:
    Schema:
      type: object
      required:
        - modelType
        - collectionName
        - displayName
        - singularName
        - pluralName
      properties:
        modelType:
          type: string
          example: collectionType
        collectionName:
          type: string
          example: users
        displayName:
          type: string
          example: Users
        singularName:
          type: string
          example: user
        pluralName:
          type: string
          example: users
        attributes:
          type: array
          items:
            $ref: '#/components/schemas/Attribute'
    Attribute:
      type: object
      required:
        - name
        - type
      properties:
        name:
          type: string
          example: username
        type:
          type: string
          example: string
        required:
          type: boolean
        max:
          type: integer
          format: int64
        min:
          type: integer
          format: int64
        maxLength:
          type: integer
          format: int32
        minLength:
          type: integer
          format: int32
        private:
          type: boolean
        notConfigurable:
          type: boolean
        default:
          type: string
        notNullable:
          type: boolean
        unsigned:
          type: boolean
        unique:
          type: boolean
    

          
          
