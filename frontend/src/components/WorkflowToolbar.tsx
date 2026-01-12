import React from 'react';

export type WorkflowStatus = 'idle' | 'running' | 'paused' | 'completed' | 'failed';

interface WorkflowToolbarProps {
  onAddNode: (type: string) => void;
  onSave: () => void;
  onTrigger?: () => void;
  onStop?: () => void;
  onResume?: () => void;
  status?: WorkflowStatus;
}

export default function WorkflowToolbar({ 
  onAddNode, 
  onSave, 
  onTrigger, 
  onStop, 
  onResume,
  status = 'idle' 
}: WorkflowToolbarProps) {
  
  const getStatusColor = (s: WorkflowStatus) => {
    switch (s) {
      case 'running': return 'bg-blue-100 text-blue-700 border-blue-200';
      case 'paused': return 'bg-yellow-100 text-yellow-700 border-yellow-200';
      case 'completed': return 'bg-green-100 text-green-700 border-green-200';
      case 'failed': return 'bg-red-100 text-red-700 border-red-200';
      default: return 'bg-gray-100 text-gray-600 border-gray-200';
    }
  };

  return (
    <div className="p-3 border-b border-gray-300 flex items-center bg-gray-50 shadow-sm gap-4">
        
        {/* Activity Palette */}
        <div className="flex items-center gap-3 overflow-x-auto max-w-xl no-scrollbar px-1 py-1 border-r border-gray-200 pr-4">
             <span className="text-gray-500 font-semibold text-sm whitespace-nowrap">Activities:</span>
             <div className="flex gap-2">
                <button 
                  className="px-3 py-1.5 bg-white border border-gray-200 rounded-md hover:bg-blue-50 hover:border-blue-200 hover:text-blue-600 text-sm font-medium transition-all shadow-sm whitespace-nowrap" 
                  onClick={() => onAddNode('Move')}
                >
                  Move
                </button>
                <button 
                  className="px-3 py-1.5 bg-white border border-gray-200 rounded-md hover:bg-blue-50 hover:border-blue-200 hover:text-blue-600 text-sm font-medium transition-all shadow-sm whitespace-nowrap" 
                  onClick={() => onAddNode('Sleep')}
                >
                  Sleep
                </button>
                <button 
                  className="px-3 py-1.5 bg-white border border-gray-200 rounded-md hover:bg-blue-50 hover:border-blue-200 hover:text-blue-600 text-sm font-medium transition-all shadow-sm whitespace-nowrap" 
                  onClick={() => onAddNode('Standup')}
                >
                  Standup
                </button>
                <button 
                  className="px-3 py-1.5 bg-white border border-gray-200 rounded-md hover:bg-blue-50 hover:border-blue-200 hover:text-blue-600 text-sm font-medium transition-all shadow-sm whitespace-nowrap" 
                  onClick={() => onAddNode('Sitdown')}
                >
                  Sitdown
                </button>
             </div>
        </div>

        <div className="flex-1"></div>

        {/* Action Buttons & Status */}
        <div className="flex gap-4 shrink-0 items-center">
            
            {/* Workflow Status Display */}
            <div className={`px-4 py-1.5 rounded-full border text-xs font-bold uppercase tracking-wider flex items-center gap-2 ${getStatusColor(status)}`}>
               <div className={`w-2 h-2 rounded-full ${status === 'running' ? 'bg-blue-500 animate-pulse' : status === 'paused' ? 'bg-yellow-500' : status === 'completed' ? 'bg-green-500' : status === 'failed' ? 'bg-red-500' : 'bg-gray-400'}`}></div>
               Status: {status}
            </div>

            {/* Separator */}
            <div className="w-px h-6 bg-gray-200"></div>

            {/* Control Group */}
            <div className="flex gap-1 border-r border-gray-200 pr-3 mr-1">
               <button 
                 onClick={onStop} 
                 disabled={status === 'idle'}
                 className="px-3 py-1.5 bg-white text-red-600 border border-gray-200 hover:bg-red-50 disabled:opacity-50 disabled:cursor-not-allowed rounded shadow-sm font-medium text-sm transition-colors flex items-center gap-1"
               >
                  <span className="w-2 h-2 bg-red-500 rounded-sm"></span> Stop
               </button>
               
               {status === 'paused' ? (
                 <button 
                   onClick={onResume} 
                   className="px-3 py-1.5 bg-white text-green-600 border border-gray-200 hover:bg-green-50 rounded shadow-sm font-medium text-sm transition-colors flex items-center gap-1"
                 >
                    <span className="w-0 h-0 border-t-[4px] border-t-transparent border-l-[6px] border-l-green-600 border-b-[4px] border-b-transparent ml-0.5"></span> Resume
                 </button>
               ) : (
                 <button 
                   onClick={onTrigger} 
                   disabled={status === 'running'}
                   className="px-3 py-1.5 bg-white text-indigo-600 border border-gray-200 hover:bg-indigo-50 disabled:opacity-50 disabled:cursor-not-allowed rounded shadow-sm font-medium text-sm transition-colors flex items-center gap-1"
                 >
                    <span className="w-0 h-0 border-t-[4px] border-t-transparent border-l-[6px] border-l-indigo-600 border-b-[4px] border-b-transparent ml-0.5"></span> Run
                 </button>
               )}
            </div>

            <button onClick={onSave} className="px-4 py-1.5 bg-green-500 text-white rounded hover:bg-green-600 shadow-sm font-medium transition-colors text-sm">Save</button>
        </div>
      </div>
  );
}
