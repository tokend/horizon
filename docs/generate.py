import yaml
import logging as log

with open('build/openapi.yaml') as f:
    data_map = yaml.safe_load(f)

schemas = data_map['components']['schemas']
new_schemas = {}


def parse_allOf(k: str, obj: dict):
    new_schema = {}

    allOf = obj.get('allOf')

    if not len(allOf):
        log.warning(f'no items in allOf, schema: {k}')
        return None, None, None

    key = allOf[0]
    if not key:
        log.warning(f'missing key in allOf, schema: {k}')
        return None, None, None

    new_schema['properties'] = {'key': key}

    if len(allOf) != 2:
        log.warning(f'unknown items in allOf, schema: {k}. Expected key and body as members of allOf field')
        return None, None, None

    schema = allOf[1]

    typ = schema.get('type')
    if typ is None:
        log.warning(f"no type in schema: {k}. Expected 'type' field in body of schema")
        return None, None, None
    new_schema['type'] = typ

    required = schema.get('required')
    if required:
        new_schema['required'] = required

    props = schema.get('properties')
    if props is None:
        log.warning(f"no properties in schema: {k}. Expected 'properties' field in body of schema")
        return None, None, None

    rels = props.get('relationships')
    if rels is None:
        log.warning(f"no relationships in schema: {k}")
    else:
        relProps = rels.get('properties')
        if relProps is None:
            log.warning(f"empty relationships in schema: {k}. Expected 'relationships' to be non-empty object")

        for relName, relValue in relProps.items():
            if 'properties' not in relValue or 'data' not in relValue.get('properties'):
                log.warning(f"no data for relationship: {relName}, schema: {k}. Expected relationships with data")
                continue

            if 'type' not in relValue['properties']['data']:
                relValue['properties']['data']['type'] = 'object'

            relProps[relName] = relValue
        rels['properties'] = relProps

        new_schema['properties']['relationships'] = {'$ref': '#/components/schemas/' + k + 'Relationships'}

    attrs = props.get('attributes')
    if attrs is None:
        log.warning(f"no attributes in schema: {k}.")
    else:
        new_schema['properties']['attributes'] = {'$ref': '#/components/schemas/' + k + 'Attributes'}

    return new_schema, rels, attrs


def get_simple_ref(raw_ref: str):
    if raw_ref.startswith('#/components/'):
        return raw_ref[1+raw_ref.rindex('/'):]


def parse_oneOf(s_name: str, obj: dict) -> dict or None:
    new_schema = {
        'type': 'object',
        'required': ['type'],
        'properties': {
            'type': {
                '$ref': '#/components/schemas/Enum'
            }
        },
    }

    oneOf = obj.get('oneOf')
    if not oneOf:
        log.warning(f'no items in allOf, schema: {s_name}')
        return new_schema

    if not len(oneOf):
        log.warning(f'no items in allOf, schema: {s_name}')
        return None

    for item in oneOf:
        if '$ref' not in item:
            log.warning(f'unknown structure of oneOf: not only $ref\' are presented: schema: {s_name}')
            return None

        simple_ref = get_simple_ref(item['$ref'])
        new_schema['properties'][simple_ref.lower()] = {
            'nullable': True,
            '$ref': item['$ref'],
        }

    return new_schema


for k, v in schemas.items():
    # if 'allOf' not in v and 'oneOf' not in v and 'properties' not in v:
    #     log.warning(f"Unknown schema structure. 'allOf' and 'properties' both are not body of the schema {k}")
    #     continue

    new_schema = {}

    allOf = v.get('allOf')
    oneOf = v.get('oneOf')
    props = v.get('properties')

    if allOf:
        new_schema, rels, attrs = parse_allOf(k, v)
        if rels:
            new_schemas[k + 'Relationships'] = rels
        if attrs:
            new_schemas[k + 'Attributes'] = attrs
    elif oneOf:
        new_schema = parse_oneOf(k, v)
    else:
        new_schema = v

    new_schemas[k] = new_schema

data_map['components']['schemas'] = new_schemas

data_map['paths'] = None
data_map['servers'] = None
data_map['info'] = None
data_map['x-tagGroups'] = None
data_map['tags'] = None

with open('build/regenerated.yaml', "w") as f:
    yaml.dump(data_map, stream=f)
