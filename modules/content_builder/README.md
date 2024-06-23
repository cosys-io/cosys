# cosys - content builder
This module is responsible for handling the creation of new content types.

### Documentation:

```yaml
paths:
  /admin/schema:
    get:
      summary: Returns all schemas
      responses:
        '200':
          description: ok
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Schema'
    post:
      summary: Create a new content type
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Schema'
        required: true
      responses:
        '200':
          description: ok
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
        - simplifiedDataType
        - detailedDataType
      properties:
        name:
          type: string
          example: age
        simplifiedDataType:
          type: string
          example: Number
        detailedDataType:
          type: string
          example: Int
        shownInTable:
          type: boolean
          default: true
        required:
          type: boolean
          default: false
        max:
          type: integer
          format: int64
          default: 2147483647
        min:
          type: integer
          format: int64
          default: -2147483648
        maxLength:
          type: integer
          format: int32
          default: -1
        minLength:
          type: integer
          format: int32
          default: -1
        private:
          type: boolean
          default: false
        editable:
          type: boolean
          default: true
        default:
          type: string
          default: ""
        nullable:
          type: boolean
          default: true
        unique:
          type: boolean
          default: false