import type { Node } from '@xyflow/react';
import type { BaseParams, MoveParams, SleepParams, HeadParams, TTSParams } from '../types/workflows';

interface NodeEditorModalProps {
  isOpen: boolean;
  node: Node | null;
  onClose: () => void;
  onSave: (params: BaseParams | MoveParams | SleepParams | HeadParams | TTSParams) => void;
  onDelete: () => void;
}

export default function NodeEditorModal({ isOpen, node, onClose, onSave, onDelete }: NodeEditorModalProps) {

  if (!isOpen || !node) return null;

  const data = node.data;

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    
    // Extract params based on type
    let params: BaseParams | MoveParams | SleepParams | HeadParams | TTSParams = {};
    const type = data.activityType;
    
    if (type === 'Move') {
      params = { 
        x: Number(formData.get('x')), 
        y: Number(formData.get('y')),
        orientation: Number(formData.get('orientation'))
      };
    } else if (type === 'Sleep') {
      params = { duration: Number(formData.get('duration')) };
    } else if (type === 'Head') {
      params = { angle: Number(formData.get('angle')) };
    } else if (type === 'TTS') {
      params = { text: String(formData.get('text')) };
    } else {
      params = {};
    }
    onSave(params);
  };

  return (
    <div className="fixed inset-0 bg-black/50 z-50 flex justify-center items-center backdrop-blur-sm">
      <div className="bg-white p-6 rounded-lg min-w-[320px] text-gray-800 shadow-2xl">
        <h3 className="mt-0 text-lg font-bold border-b pb-2 mb-4">Edit {data.activityType}</h3>
        
        <form onSubmit={handleSubmit}>
        
        {/* Move Params */}
        {data.activityType === 'Move' && (
          <>
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">X Coordinate:</label>
              <input name="x" type="number" step="0.1" defaultValue={(data.params as MoveParams).x} className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 outline-none" />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">Y Coordinate:</label>
              <input name="y" type="number" step="0.1" defaultValue={(data.params as MoveParams).y} className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 outline-none" />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">Orientation(degree):</label>
              <input name="orientation" type="number" step="0.1" defaultValue={(data.params as MoveParams).orientation} className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 outline-none" />
            </div>
          </>
        )}

         {/* Sleep Params */}
         {data.activityType === 'Sleep' && (
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">Duration (ms):</label>
              <input name="duration" type="number" defaultValue={(data.params as SleepParams).duration} className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 outline-none" />
            </div>
        )}

        {/* Head Params*/}
        {data.activityType === 'Head' && (
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">Angle (degrees):</label>
              <input name="angle" type="number" defaultValue={(data.params as HeadParams).angle} className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 outline-none" />
            </div>
        )}

        {/* TTS Params*/}
        {data.activityType === 'TTS' && (
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">Text:</label>
              <input name="text" type="text" defaultValue={(data.params as TTSParams).text} className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 outline-none" />
            </div>
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
