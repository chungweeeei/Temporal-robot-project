import { Handle, Position, type NodeProps } from '@xyflow/react';
import type { FlowNodeData } from '../../types/schema';

export default function StartNode({ data }: NodeProps) {
  const { label } = data as FlowNodeData;
  return (
    <div className="flex flex-col items-center justify-center">
        <div className="w-16 h-16 rounded-full bg-gradient-to-br from-green-400 to-green-600 shadow-lg border-4 border-white flex items-center justify-center hover:scale-110 transition-transform">
             <div className="text-white font-bold text-xs uppercase tracking-wider">{label}</div>
        </div>
      
      {/* Output Handle */}
      <Handle 
        type="source" 
        position={Position.Right} 
        className="!w-4 !h-4 !bg-green-500 !border-2 !border-white !-right-2"
      />
    </div>
  );
}
