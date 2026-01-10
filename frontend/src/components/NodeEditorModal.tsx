import React, { useState, useEffect } from 'react';
import type { Node } from '@xyflow/react';
import type { FlowNodeData, BaseParams, RetryPolicy, ConditionParams, MoveParams, SleepParams } from '../types/schema';

interface NodeEditorModalProps {
  isOpen: boolean;
  node: Node | null;
  onClose: () => void;
  onSave: (params: BaseParams | MoveParams | SleepParams | ConditionParams, retryPolicy: RetryPolicy) => void;
  onDelete: () => void;
}

export default function NodeEditorModal({ isOpen, node, onClose, onSave, onDelete }: NodeEditorModalProps) {
  const [retryEnabled, setRetryEnabled] = useState(false);
  const [localMaxAttempts, setLocalMaxAttempts] = useState(1);

  // Reset state when node changes
  useEffect(() => {
    if (node && isOpen) {
        const data = node.data as FlowNodeData;
        const policy = data.retryPolicy;
        const currentAttempts = policy?.maxAttempts ?? 1;
        
        setLocalMaxAttempts(currentAttempts);
        setRetryEnabled(currentAttempts > 1);
    }
  }, [node, isOpen]);

  if (!isOpen || !node) return null;

  const data = node.data as FlowNodeData;

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    
    // Extract params based on type
    let params: BaseParams | MoveParams | SleepParams | ConditionParams = {};
    const type = data.activityType;
    
    if (type === 'Move') {
      params = { 
        x: Number(formData.get('x')), 
        y: Number(formData.get('y')) 
      };
    } else if (type === 'Sleep') {
      params = { duration: Number(formData.get('duration')) };
    } else if (type === 'Condition') {
      params = { expression: formData.get('expression') as string };
    } else {
      params = {};
    }

    // Updated Retry Policy Logic
    const retryPolicy: RetryPolicy = {
      maxAttempts: retryEnabled ? localMaxAttempts : 1,
      initialInterval: (retryEnabled && localMaxAttempts !== 0) ? Number(formData.get('initialInterval')) : 1000,
    };

    onSave(params, retryPolicy);
  };

  return (
    <div className="fixed inset-0 bg-black/50 z-50 flex justify-center items-center backdrop-blur-sm">
      <div className="bg-white p-6 rounded-lg min-w-[320px] text-gray-800 shadow-2xl">
        <h3 className="mt-0 text-lg font-bold border-b pb-2 mb-4">Edit {data.activityType}</h3>
        
        <form onSubmit={handleSubmit}>

        {/* Condition Params */}
        {data.activityType === 'Condition' && (
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">Expression (e.g. x {'>'} 5):</label>
              <input name="expression" type="text" defaultValue={(data.params as ConditionParams).expression} className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 outline-none" />
            </div>
        )}
        
        {/* Move Params */}
        {data.activityType === 'Move' && (
          <>
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">X Coordinate:</label>
              <input name="x" type="number" defaultValue={(data.params as MoveParams).x} className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 outline-none" />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">Y Coordinate:</label>
              <input name="y" type="number" defaultValue={(data.params as MoveParams).y} className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 outline-none" />
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

        {/* Common Retry Policy Settings */}
        <div className="mt-6 pt-4 border-t border-gray-100">
            <div className="flex items-center justify-between mb-3">
              <strong className="text-sm text-gray-600">Retry Policy</strong>
              <label className="flex items-center gap-2 cursor-pointer select-none">
                <div className={`w-9 h-5 rounded-full p-0.5 transition-colors ${retryEnabled ? 'bg-blue-500' : 'bg-gray-300'}`}>
                  <div className={`w-4 h-4 rounded-full bg-white shadow-sm transform transition-transform ${retryEnabled ? 'translate-x-4' : 'translate-x-0'}`} />
                </div>
                <input 
                  type="checkbox" 
                  className="hidden" 
                  checked={retryEnabled} 
                  onChange={(e) => setRetryEnabled(e.target.checked)} 
                />
                <span className="text-xs text-gray-500 font-medium">{retryEnabled ? 'On' : 'Off'}</span>
              </label>
            </div>

            {retryEnabled && (
              <div className="flex gap-4 p-4 bg-gray-50 rounded-lg border border-gray-100 animate-in fade-in slide-in-from-top-2 duration-200">
                <div className="flex-1">
                   <label className="block text-xs mb-1.5 text-gray-500 font-medium">Max Attempts (0 = Infinite)</label>
                   <input 
                      name="maxAttempts" 
                      type="number" 
                      value={localMaxAttempts} 
                      onChange={(e) => setLocalMaxAttempts(parseInt(e.target.value) || 0)}
                      className="w-full p-2 border border-gray-300 rounded text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all" 
                      min="0" 
                   />
                </div>
                <div className="flex-1">
                   <label className={`block text-xs mb-1.5 font-medium ${localMaxAttempts === 0 ? 'text-gray-400' : 'text-gray-500'}`}>
                     Initial Interval (ms)
                   </label>
                   <input 
                      name="initialInterval" 
                      type="number" 
                      defaultValue={(data.retryPolicy as RetryPolicy)?.initialInterval || 1000} 
                      disabled={localMaxAttempts === 0}
                      className={`w-full p-2 border rounded text-sm outline-none transition-all ${
                        localMaxAttempts === 0 
                          ? 'bg-gray-100 border-gray-200 text-gray-400 cursor-not-allowed' 
                          : 'bg-white border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-blue-500'
                      }`} 
                   />
                </div>
              </div>
            )}
        </div>

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
