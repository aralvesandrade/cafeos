import { Sprout, DollarSign, MapPin, Grid3X3, TrendingUp } from "lucide-react";

interface Stats {
  totalFarms: number;
  totalPlots: number;
  totalProduction: number;
  totalCost: number;
  coe?: number;
  cot?: number;
  ctProducao?: number;
}

export function StatsCards({ stats }: { stats: Stats }) {
  const items = [
    { icon: MapPin, label: "Fazendas", value: stats.totalFarms },
    { icon: Grid3X3, label: "Talhões", value: stats.totalPlots },
    {
      icon: Sprout,
      label: "Produção (sacas)",
      value: stats.totalProduction.toLocaleString(),
    },
    {
      icon: DollarSign,
      label: "Custo total",
      value: `R$ ${stats.totalCost.toLocaleString()}`,
    },
    ...(stats.coe
      ? [
          {
            icon: TrendingUp,
            label: "COE (Custo Operacional Efetivo)",
            value: `R$ ${stats.coe.toLocaleString()}`,
            title:
              "Custo Operacional Efetivo: soma dos desembolsos operacionais diretos (insumos, mão de obra contratada, combustível etc.)",
          },
        ]
      : []),
    ...(stats.cot
      ? [
          {
            icon: TrendingUp,
            label: "COT (Custo Operacional Total)",
            value: `R$ ${stats.cot.toLocaleString()}`,
            title:
              "Custo Operacional Total: COE + mão de obra familiar + depreciação de máquinas e benfeitorias",
          },
        ]
      : []),
    ...(stats.ctProducao
      ? [
          {
            icon: TrendingUp,
            label: "CT (Custo Total)",
            value: `R$ ${stats.ctProducao.toLocaleString()}`,
            title:
              "Custo Total: COT + remuneração do capital investido (terra e outros ativos)",
          },
        ]
      : []),
  ];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {items.map((item) => (
        <div
          key={item.label}
          title={item.title}
          className="bg-card border border-border rounded-lg px-4 pt-3 pb-4 border-t-2 border-t-primary"
        >
          <div className="flex items-center justify-between mb-2">
            <p className="text-[11px] uppercase tracking-widest text-muted-foreground">{item.label}</p>
            <item.icon className="h-3.5 w-3.5 text-muted-foreground" />
          </div>
          <p className="font-mono text-2xl font-medium text-foreground tabular-nums">{item.value}</p>
        </div>
      ))}
    </div>
  );
}
