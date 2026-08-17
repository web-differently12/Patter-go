import React from 'react';

export const CallLogsPage: React.FC = () => {
  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      <div className="border-b border-zinc-800 pb-4">
        <span className="text-xs font-mono font-bold text-cyan-400">Call History & Audio Inspector</span>
        <h1 className="text-2xl font-bold tracking-tight text-white mt-0.5">Historique des Appels & Lecteur Audio</h1>
      </div>

      <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 shadow-xl space-y-4 text-xs">
        <div className="flex justify-between items-center border-b border-zinc-800/80 pb-3">
          <span className="font-mono text-zinc-300">Appels récents (Twilio / Telnyx / Plivo)</span>
          <span className="text-[10px] text-emerald-400 font-mono">Realtime Stream Active</span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-zinc-300 font-mono text-[11px]">
            <thead className="bg-zinc-900 text-zinc-400 border-b border-zinc-800">
              <tr>
                <th className="p-3">Call ID</th>
                <th className="p-3">Direction</th>
                <th className="p-3">Numéro</th>
                <th className="p-3">Durée</th>
                <th className="p-3 text-right">Coût</th>
                <th className="p-3 text-center">Audio WAV</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/60">
              <tr className="hover:bg-zinc-900/40">
                <td className="p-3 text-cyan-400">call_98234_abc</td>
                <td className="p-3 text-emerald-400 font-bold">INBOUND</td>
                <td className="p-3">+33 6 12 34 56 78</td>
                <td className="p-3">02:14</td>
                <td className="p-3 text-right">0.082 €</td>
                <td className="p-3 text-center">
                  <button className="px-2.5 py-1 rounded bg-zinc-900 border border-zinc-800 hover:border-cyan-500 text-zinc-300">▶ Ecouter</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
