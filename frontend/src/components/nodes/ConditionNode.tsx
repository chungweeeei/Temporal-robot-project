import React from 'react';
import { Handle, Position, type NodeProps } from '@xyflow/react';
import type { FlowNodeData } from '../../types/schema';

export default function ConditionNode({ data }: NodeProps) {
  const { label } = data as FlowNodeData;
  return (
    <div className="flex bg-white rounded-lg shadow-md border-2 border-amber-400 w-[220px] hover:shadow-xl hover:scale-105 transition-all">
      {/* Input Handle (Left) */}
      <Handle type="target" position={Position.Left} className="!w-3 !h-3 !bg-gray-400 !border-2 !border-white !-left-1.5" />
      
      {/* Left: Content */}
      <div className="flex-1 p-3 flex flex-col justify-center border-r border-amber-100 bg-amber-50 rounded-l-md">
         <div className="flex items-center gap-1.5 mb-1.5">
            <div className="w-2 h-2 bg-amber-500 rotate-45 shrink-0"></div>
            <div className="font-bold text-amber-900 text-sm truncate w-full" title={label}>{label}</div>
         </div>
         <code className="text-[10px] font-mono bg-white px-2 py-1 rounded border border-amber-200 text-gray-600 block truncate">
            {(data.params as any)?.expression || 'EXP'}
         </code>
      </div>

      {/* Right: Outputs */}
      <div className="w-[70px] flex flex-col justify-between py-2 bg-white rounded-r-lg">
           {/* True (Top Right) */}
           <div className="flex items-center justify-end px-2 relative">
              <span className="text-[9px] font-bold text-green-600 uppercase tracking-wide">True</span>
              <Handle 
                type="source" 
                position={Position.Right} 
                id="true" 
                className="!w-3 !h-3 !bg-green-500 !border-2 !border-white !-right-1.5"
              />
           </div>

           {/* False (Bottom Right) */}
           <div className="flex items-center justify-end px-2 relative">
              <span className="text-[9px] font-bold text-red-600 uppercase tracking-wide">False</span>
              <Handle 
                type="source" 
                position={Position.Right} 
                id="false" 
                className="!w-3 !h-3 !bg-red-500 !border-2 !border-white !-right-1.5"
              />
           </div>
      </div>
    </div>
  );
}
