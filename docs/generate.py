import yaml
import logging as log
import argparse

parser = argparse.ArgumentParser(description='Generate intermediate .yaml for regources generator')
parser.add_argument('--path', type=str, help='path to generated openapi.yaml from horizon docs')
parser.add_argument('--out', type=str, help='path to write generated intermediare .yaml file')

args = parser.parse_args()

with open(args.path) as f:
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


def get_referenced_schema(simple_ref: str):
    for k, v in schemas.items():
        if k == simple_ref:
            return v
    return None


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
            log.warning(f'unknown structure of oneOf: $ref\' is not presented: schema: {s_name}')
            return None

        simple_ref = get_simple_ref(item['$ref'])
        new_schema['properties'][simple_ref.lower()] = {
            'nullable': True,
            '$ref': item['$ref'],
        }

        typ = item.get('type')
        fmt = item.get('format')
        props = item.get('properties')
        if props:
            if 'type' in props:
                if 'format' in props.get('type'):
                    fmt = props.get('type').get('format')

        if typ:
            new_schema['properties'][simple_ref.lower()]['type'] = typ
        if fmt:
            new_schema['properties'][simple_ref.lower()]['format'] = fmt

    return new_schema


for k, v in schemas.items():
    if 'allOf' not in v and 'oneOf' not in v and 'properties' not in v:
        log.warning(f"Unknown schema structure. 'allOf' and 'properties' both are not body of the schema {k}")
        continue

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

with open(args.out, "w") as f:
    yaml.dump(data_map, stream=f)
