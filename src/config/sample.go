package config

import (
	"alc/model"
)

var (
	VidrioCategories []model.Category = []model.Category{
		{
			Id:   1,
			Type: model.VidrioType,
			Name: "Vidrio Monolítico",
			Description: `Lámina de vidrio, que puede ser incoloro, de color o
extraclaro, y se usa en espesores de 2 a 19 mm.`,
			LongDescription: `El Cristal Float se utiliza en una gran variedad de productos y proyectos, para Fachadas
e Interiores, con sistemas monolíticos. Como herramienta de diseño sus posibilidades están
sólo limitadas por la creatividad del usuario y por los criterios de seguridad, que siempre
deben ser tomados en cuenta en todas las aplicaciones del cristal plano para arquitectura y
componentes de equipamiento.`,
			Slug: "vidrio-monolitico",
		},
		{
			Id:   2,
			Type: model.VidrioType,
			Name: "Vidrio Reflectivo",
			Description: `Ofrece control solar, ahorro energético y privacidad
con un efecto espejado que varía según la luz del día.`,
			Slug: "vidrio-reflectivo",
		},
		{
			Id:   3,
			Type: model.VidrioType,
			Name: "Vidrio Decorativo",
			Description: `Para proyectos de multiples diseños, incluye técnicas
como arenado, esmerilado, serigrafiado y mateado al ácido.`,
			Slug: "vidrio-decorativo",
		},
	}

	MonoliticoFeatures []model.CategoryFeature = []model.CategoryFeature{
		{
			Id:          1,
			Category:    VidrioCategories[0],
			Name:        "Tipo de Vidrio",
			Description: "Vidrio monolítico (vidrio float o recocido)",
		},
		{
			Id:          1,
			Category:    VidrioCategories[0],
			Name:        "Resistencia Mecánica",
			Description: "Resistencia a flexión: 20-30 MPa",
		},
		{
			Id:          1,
			Category:    VidrioCategories[0],
			Name:        "Resistencia Térmica",
			Description: "Resistencia a choque térmico: 40-50°C",
		},
	}

	MonoliticoItems []model.Item = []model.Item{
		{
			Id:          1,
			Category:    VidrioCategories[0],
			Name:        "2 - 3 mm",
			Description: "Espejos decorativos, muebles con vidrio, cuadros y marcos.",
		},
		{
			Id:          2,
			Category:    VidrioCategories[0],
			Name:        "4 - 5 mm",
			Description: "Ventanas pequeñas, divisiones interiores, mamparas de baño, vitrinas.",
		},
		{
			Id:          3,
			Category:    VidrioCategories[0],
			Name:        "6 - 8 mm",
			Description: "Ventanas medianas, puertas de vidrio, barras de cocina, escalones.",
		},
	}

	AluminioCategories []model.Category = []model.Category{
		{
			Id:          1,
			Type:        model.AluminioType,
			Name:        "Sistema de Fachadas ligeras",
			Description: "",
			Slug:        "sistema-de-fachadas-ligeras",
		},
		{
			Id:          2,
			Type:        model.AluminioType,
			Name:        "Sistema de Carpintería Tecnica",
			Description: "",
			Slug:        "sistema-de-carpinteria-tecnica",
		},
		{
			Id:          3,
			Type:        model.AluminioType,
			Name:        "Sistema de Carpintería Perimetral Europea",
			Description: "",
			Slug:        "sistema-de-carpinteria-perimetral-europea",
		},
	}

	FachadasItems []model.Item = []model.Item{
		{
			Id:       1,
			Category: AluminioCategories[0],
			Name:     "Sistema Estándar",
			LongDescription: `Estas fachadas utilizan sistemas estándar de aluminio disponibles comercialmente, diseñados
para edificaciones pequeñas o medianas sin requerimientos constructivos complejos.
Son ligeras, con montantes y travesaños, y adecuadas para edificios de hasta 10 pisos.
Como muro cortina, se instalan frente a los forjados y quedan suspendidas ("Glass Skin").
Pueden usar vidrios laminados para control acústico o templados de hasta 10 mm en diversos
colores.`,
			Slug: "sistema-estandar",
		},
		{
			Id:       2,
			Category: AluminioCategories[0],
			Name:     "Sistema Stick",
			Slug:     "sistema-stick",
		},
		{
			Id:       3,
			Category: AluminioCategories[0],
			Name:     "Sistema Frame",
			Slug:     "sistema-frame",
		},
	}

	UPVCCategories []model.Category = []model.Category{
		{
			Id:          1,
			Type:        model.UPVCType,
			Name:        "Lumina 60 Ventanas Articuladas",
			Description: "",
			Slug:        "lumina-60-ventanas-articuladas",
		},
		{
			Id:          2,
			Type:        model.UPVCType,
			Name:        "Lumina 104 Puertas",
			Description: "",
			Slug:        "lumina-104-puertas",
		},
		{
			Id:          3,
			Type:        model.UPVCType,
			Name:        "Natura 66 Ventanas Correderas",
			Description: "",
			Slug:        "natura-66-ventanas-correderas",
		},
	}

	LuminaItems []model.Item = []model.Item{
		{
			Id:       1,
			Category: UPVCCategories[0],
			Name:     "Proyectante",
			LongDescription: `La Ventana Proyectante de Perfil uPVC es ideal para oficinas y edificaciones que requieren alta
eficiencia térmica y acústica. Su apertura hacia el exterior mediante brazos de extensión
permite una ventilación controlada sin corrientes de aire y evita la entrada de agua en caso
de lluvia.`,
			Slug: "proyectante",
		},
		{
			Id:       2,
			Category: UPVCCategories[0],
			Name:     "Batiente",
			Slug:     "batiente",
		},
		{
			Id:       3,
			Category: UPVCCategories[0],
			Name:     "Oscilobatiente",
			Slug:     "oscilobatiente",
		},
	}

	Imgs []model.Image = []model.Image{
		{
			Id: 0,
		},
		{
			Id: 0,
		},
		{
			Id: 0,
		},
	}
)
