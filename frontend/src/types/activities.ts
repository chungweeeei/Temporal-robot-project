export interface SchemaProperty {
    type: 'string' | 'number' | 'boolean' | 'object' | 'array';
    title?: string;
    default?: any;
}

export interface InputSchema {
    type: "object";
    required?: string[];
    properties: Record<string, SchemaProperty>
}

export interface ActivityDefinition {
    id: number;
    name: string;
    activity_type: string;
    node_type: string;
    input_schema: InputSchema | null;  // 可能為 null（如 Standup, Sitdown）
    created_at: string;
    updated_at: string;
}