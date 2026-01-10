import React from 'react';
import { Handle, Position, type NodeProps } from '@xyflow/react';
import type { FlowNodeData } from '../../types/schema';

export default function ActionNode({ data }: NodeProps) {
  const { label } = data as FlowNodeData;
  
  return (
    <div className="flex bg-white rounded-xl shadow-md border border-gray-200 w-[180px] transition-all hover:shadow-xl hover:scale-105 hover:border-blue-400">
      {/* Input Handle (Left) */}
      <Handle type="target" position={Position.Left} className="!w-3 !h-3 !bg-gray-400 !border-2 !border-white !-left-1.5" />
      
      {/* Left: Content */}
      <div className="flex-1 p-3 flex flex-col justify-center border-r border-gray-100">
         <div className="font-bold text-gray-800 text-sm truncate mb-1" title={label}>{label}</div>
         <div className="text-[10px] text-gray-500 uppercase font-bold tracking-wider">Action</div>
      </div>

      {/* Right: Outputs */}
      <div className="w-[80px] flex flex-col justify-between py-2 bg-gray-50 rounded-r-xl">
          {/* Success (Top Right) */}
           <div className="flex items-center justify-end px-3 relative">
              <span className="text-[9px] font-bold text-green-600 uppercase tracking-wide">Success</span>
              <Handle 
                type="source" 
                position={Position.Right} 
                id="success" 
                className="!w-3 !h-3 !bg-green-500 !border-2 !border-white !-right-1.5 transition-all hover:!w-4 hover:!h-4"
              />
           </div>

           {/* Failure (Bottom Right) */}
           <div className="flex items-center justify-end px-3 relative">
              <span className="text-[9px] font-bold text-red-600 uppercase tracking-wide">Failure</span>
              <Handle 
                type="source" 
                position={Position.Right} 
                id="failure" 
                className="!w-3 !h-3 !bg-red-500 !border-2 !border-white !-right-1.5 transition-all hover:!w-4 hover:!h-4"
              />
           </div>
      </div>
    </div>
  );
}
