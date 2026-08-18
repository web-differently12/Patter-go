import React, { useState } from 'react';

export const PricingPublicPage: React.FC = () => {
  const [billingCycle, setBillingCycle] = useState<'monthly' | 'yearly'>('monthly');
  const [estimatedMinutes, setEstimatedMinutes] = useState<number>(2500);

  // 1 EUR = 1000 Credits. Voice Standard = 15 Credits/min.
  const requiredCredits = estimatedMinutes * 15;
  const equivalentEUR = requiredCredits / 1000;

  return (
    <div className="min-h-screen bg-[#09090b] text-zinc-100 font-sans p-6 space-y-12 max-w-7xl mx-auto">
      {/* Header */}
      <div className="text-center space-y-3 pt-8">
        <span className="text-xs font-mono font-bold text-cyan-400 bg-zinc-950 px-3 py-1 rounded-full border border-zinc-800">
          Tarification Transparente
        </span>
        <h1 className="text-3xl sm:text-5xl font-extrabold tracking-tight text-zinc-100">
          Plans Stripe et Bareme de Credits Unifie
        </h1>
        <p className="text-xs sm:text-sm text-zinc-400 max-w-xl mx-auto">
          Choisissez un plan d'abonnement Mensuel ou Annuel. 1 EUR = 1 000 Credits KallFlow.
        </p>

        {/* Monthly / Yearly Switch */}
        <div className="pt-4 flex justify-center">
          <div className="bg-zinc-950 border border-zinc-800 p-1 rounded-xl flex items-center gap-1 text-xs">
            <button
              onClick={() => setBillingCycle('monthly')}
              className={`px-4 py-1.5 rounded-lg font-semibold transition-all ${
                billingCycle === 'monthly' ? 'bg-zinc-100 text-zinc-950 font-bold' : 'text-zinc-400 hover:text-white'
              }`}
            >
              Facturation Mensuelle
            </button>
            <button
              onClick={() => setBillingCycle('yearly')}
              className={`px-4 py-1.5 rounded-lg font-semibold transition-all ${
                billingCycle === 'yearly' ? 'bg-cyan-500 text-zinc-950 font-bold' : 'text-zinc-400 hover:text-white'
              }`}
            >
              Facturation Annuelle (-20%)
            </button>
          </div>
        </div>
      </div>

      {/* Pricing Cards Grid */}
      <div id="pricing" className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Card 1: Starter */}
        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-6 shadow-xl space-y-4 flex flex-col justify-between">
          <div className="space-y-3">
            <h3 className="font-bold text-white text-base">Starter / Dev</h3>
            <p className="text-xs text-zinc-400">Pour developper et tester des agents vocaux et messaging.</p>
            <div className="text-3xl font-extrabold font-mono text-zinc-100">
              {billingCycle === 'monthly' ? '49 EUR' : '470 EUR'} <span className="text-xs text-zinc-500 font-sans font-normal">/ {billingCycle === 'monthly' ? 'mois' : 'an'}</span>
            </div>
            <ul className="text-xs text-zinc-300 space-y-2 pt-2 border-t border-zinc-800/80">
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> 50 000 Credits mensuels inclus</li>
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> 2 500 Messages messaging gratuits</li>
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> 1 Instance messaging active</li>
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> 1 Client MCP / API externe</li>
            </ul>
          </div>
          <button className="w-full py-2.5 rounded-xl font-bold bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 text-zinc-200 text-xs transition-colors">
            Souscrire Starter →
          </button>
        </div>

        {/* Card 2: Agency Pro (Highlighted) */}
        <div className="bg-zinc-950 border border-cyan-500/80 rounded-2xl p-6 shadow-2xl space-y-4 flex flex-col justify-between relative overflow-hidden">
          <span className="absolute top-3 right-3 text-[10px] font-bold px-2 py-0.5 rounded bg-cyan-500 text-zinc-950">Best-Seller Agence</span>
          <div className="space-y-3">
            <h3 className="font-bold text-white text-base">Agency White-Label</h3>
            <p className="text-xs text-zinc-400">Pour les agences et resellers en Marque Blanche.</p>
            <div className="text-3xl font-extrabold font-mono text-zinc-100">
              {billingCycle === 'monthly' ? '149 EUR' : '1 430 EUR'} <span className="text-xs text-zinc-500 font-sans font-normal">/ {billingCycle === 'monthly' ? 'mois' : 'an'}</span>
            </div>
            <ul className="text-xs text-zinc-300 space-y-2 pt-2 border-t border-zinc-800/80">
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> 150 000 Credits mensuels inclus</li>
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> 15 000 Messages messaging gratuits</li>
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> 5 Instances messaging actives</li>
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> 5 Clients MCP / API externes</li>
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> Marque Blanche & Marge personnalisable</li>
            </ul>
          </div>
          <button className="w-full py-2.5 rounded-xl font-bold bg-cyan-500 hover:bg-cyan-400 text-zinc-950 text-xs transition-colors shadow-lg">
            Activer Plan Agency →
          </button>
        </div>

        {/* Card 3: Enterprise */}
        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-6 shadow-xl space-y-4 flex flex-col justify-between">
          <div className="space-y-3">
            <h3 className="font-bold text-white text-base">Enterprise Scale</h3>
            <p className="text-xs text-zinc-400">Pour les grands comptes et fort volume d'appels.</p>
            <div className="text-3xl font-extrabold font-mono text-zinc-100">
              {billingCycle === 'monthly' ? '499 EUR' : '4 790 EUR'} <span className="text-xs text-zinc-500 font-sans font-normal">/ {billingCycle === 'monthly' ? 'mois' : 'an'}</span>
            </div>
            <ul className="text-xs text-zinc-300 space-y-2 pt-2 border-t border-zinc-800/80">
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> 500 000 Credits mensuels inclus</li>
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> Messages messaging Illimites</li>
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> 25 Instances messaging actives</li>
              <li className="flex items-center gap-2"><span className="text-cyan-400">✓</span> Clients MCP Illimites</li>
            </ul>
          </div>
          <button className="w-full py-2.5 rounded-xl font-bold bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 text-zinc-200 text-xs transition-colors">
            Contacter Enterprise →
          </button>
        </div>
      </div>

      {/* Interactive Credits Estimator */}
      <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-6 shadow-xl space-y-4">
        <h3 className="font-bold text-white text-base">Simulateur de Consommation en Credits KallFlow</h3>
        <p className="text-xs text-zinc-400">Estimez vos besoins mensuels en minutes de conversation vocale standard (15 Credits / min).</p>

        <div className="space-y-3 pt-2">
          <div className="flex justify-between items-center text-xs font-mono">
            <span>Volume souhaite: <strong>{estimatedMinutes.toLocaleString()} minutes / mois</strong></span>
            <span className="text-cyan-400 font-bold">{requiredCredits.toLocaleString()} Credits ({equivalentEUR.toFixed(2)} EUR)</span>
          </div>
          <input
            type="range"
            min={100}
            max={20000}
            step={100}
            value={estimatedMinutes}
            onChange={(e) => setEstimatedMinutes(Number(e.target.value))}
            className="w-full h-2 bg-zinc-800 rounded-lg appearance-none cursor-pointer accent-cyan-500"
          />
        </div>
      </div>
    </div>
  );
};
