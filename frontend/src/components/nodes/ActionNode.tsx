import { Handle, Position, type NodeProps } from '@xyflow/react';

export default function ActionNode({ data }: NodeProps) {
  const label = data.label as string;
  const isRunning = data.isRunning as boolean;
  
  return (
    <div 
      className={`
        flex bg-white rounded-xl shadow-lg border-2 w-[180px] transition-all
        ${isRunning 
          ? 'border-blue-500 shadow-blue-300 shadow-xl ring-2 ring-blue-400 ring-offset-2 animate-pulse' 
          : 'border-slate-300 hover:shadow-xl hover:scale-105 hover:border-blue-400 shadow-slate-200'
        }
      `}
    >
      {/* Input Handle (Left) */}
      <Handle 
        type="target" 
        position={Position.Left} 
        className={`!w-3 !h-3 !border-2 !border-white !-left-1.5 ${isRunning ? '!bg-blue-500' : '!bg-gray-400'}`} 
      />
      
      {/* Left: Content */}
      <div className="flex-1 p-3 flex flex-col justify-center border-r border-gray-100">
        <div className="flex items-center gap-2">
          {isRunning && (
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-blue-500"></span>
            </span>
          )}
          <div className="font-bold text-gray-800 text-sm truncate" title={label}>{label}</div>
        </div>
        <div className={`text-[10px] uppercase font-bold tracking-wider ${isRunning ? 'text-blue-500' : 'text-gray-500'}`}>
          {isRunning ? 'Running...' : 'Action'}
        </div>
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
