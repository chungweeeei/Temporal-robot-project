import { useEffect, useState } from 'react';
import type { Node } from '@xyflow/react';
import type { InputSchema } from '../types/activities';
import type { BaseParams, MoveParams, SleepParams, HeadParams, TTSParams } from '../types/workflows';

interface NodeEditorModalProps {
  isOpen: boolean;
  node: Node | null;
  onClose: () => void;
  onSave: (params: BaseParams | MoveParams | SleepParams | HeadParams | TTSParams) => void;
  onDelete: () => void;
}

export default function NodeEditorModal({ isOpen, node, onClose, onSave, onDelete }: NodeEditorModalProps) {

  const [formValues, setFormValues] = useState<Record<string, any>>({});

  const data = node?.data as {
    label: string;
    activityType: string;
    inputSchema: InputSchema | null;
    params: BaseParams | MoveParams | SleepParams | HeadParams | TTSParams;
  };

  useEffect(() => {
    if (data?.params) {
      setFormValues(data.params);
    }
  }, [node, data?.params]);

  if (!isOpen || !node) return null;

  const { inputSchema } = data; 

  const handleChange = (key: string, value: any, type: string) => {
    setFormValues(prev => ({
      ...prev,
      [key]: type === 'number' ? Number(value) : value
    }));
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    onSave(formValues);
  };

  // Render form fields based on input schema
  const renderField = (key: string, property: { type: string; title?: string; default?: any }) => {
    const value = formValues[key] ?? property.default ?? '';
    const label = property.title || key;
    const isRequired = inputSchema?.required?.includes(key);

    return (
      <div key={key} className="mb-4">
        <label className="block text-sm font-medium mb-1">
          {label}
          {isRequired && <span className="text-red-500 ml-1">*</span>}
        </label>
        {property.type === 'string' ? (
          <input
            name={key}
            type="text"
            value={value}
            onChange={(e) => handleChange(key, e.target.value, property.type)}
            required={isRequired}
            className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 outline-none"
          />
        ) : property.type === 'number' ? (
          <input
            name={key}
            type="number"
            step="any"
            value={value}
            onChange={(e) => handleChange(key, e.target.value, property.type)}
            required={isRequired}
            className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 outline-none"
          />
        ) : property.type === 'boolean' ? (
          <input
            name={key}
            type="checkbox"
            checked={!!value}
            onChange={(e) => handleChange(key, e.target.checked, property.type)}
            className="w-5 h-5"
          />
        ) :(
          <input
            name={key}
            type="text"
            value={value}
            onChange={(e) => handleChange(key, e.target.value, property.type)}
            className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 outline-none"
          />
        )}
      </div>
    )
  };

  return (
    <div className="fixed inset-0 bg-black/50 z-50 flex justify-center items-center backdrop-blur-sm">
      <div className="bg-white p-6 rounded-lg min-w-[320px] text-gray-800 shadow-2xl">
        <h3 className="mt-0 text-lg font-bold border-b pb-2 mb-4">
          Edit {data.activityType}
        </h3>
        
        <form onSubmit={handleSubmit}>
          {inputSchema?.properties ? (
            Object.entries(inputSchema.properties).map(([key, property]) =>
              renderField(key, property)
            )
          ) : (
            <p className="text-gray-500 text-sm mb-4">
              This action has no configurable parameters.
            </p>
          )}
          <div className="flex justify-end gap-3 mt-6">
             <button type="button" onClick={onDelete} className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 transition-colors text-sm">Delete</button>
             <button type="button" onClick={onClose} className="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded transition-colors text-sm">Cancel</button>
             <button type="submit" className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors text-sm">Save</button>
          </div>
        </form>

      </div>
    </div>
  );
}
