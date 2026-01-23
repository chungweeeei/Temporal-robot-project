import { Handle, Position, type NodeProps } from '@xyflow/react';

export default function EndNode({ data }: NodeProps) {
  const label = data.label as string;
  
  return (
    <div className="flex flex-col items-center justify-center">
      {/* Input Handle */}
       <Handle 
        type="target" 
        position={Position.Left} 
        className="!w-4 !h-4 !bg-gray-800 !border-2 !border-white !-left-2"
      />

        <div className="w-16 h-16 rounded-full bg-gradient-to-br from-gray-700 to-gray-900 shadow-lg border-4 border-white flex items-center justify-center hover:scale-110 transition-transform">
             <div className="text-white font-bold text-xs uppercase tracking-wider">{label}</div>
             {/* Inner circle or icon could go here */}
             <div className="absolute w-12 h-12 rounded-full border-2 border-dashed border-gray-500 opacity-30"></div>
        </div>
    </div>
  );
}
