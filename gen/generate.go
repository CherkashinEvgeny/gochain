package main

import (
	"go/types"

	. "github.com/CherkashinEvgeny/gogen"
	tgen "github.com/CherkashinEvgeny/gogen/types"
	"github.com/CherkashinEvgeny/gogen/utils"
)

const copyright = "Code generated by github.com/CherkashinEvgeny/gochain. DO NOT EDIT."

type config struct {
	DstPkgName     string
	DstPackagePath string
	SrcPkg         *types.Package
	ImportSrcPkg   bool
	Chains         []chainConfig
}

type chainConfig struct {
	IfaceName string
	Iface     *types.Interface
	ChainName string
}

func generate(cfg config) (code string, err error) {
	imports := Imports()
	if cfg.DstPackagePath != "" {
		imports.Add(SmartImport("", "", cfg.DstPackagePath))
	}
	chains := Blocks()
	for _, chainCfg := range cfg.Chains {
		chains.Add(generateInterfaceChain(cfg.SrcPkg, chainCfg))
	}
	pkg := Pkg(copyright, cfg.DstPkgName, imports, chains)
	return Stringify(pkg), nil
}

func generateInterfaceChain(pkg *types.Package, chainCfg chainConfig) Code {
	elemName := utils.Private(chainCfg.ChainName + "Elem")
	return Blocks(
		Type(chainCfg.ChainName, Struct(FieldDecls(
			FieldDecl("root", Id(elemName)),
		))),
		generateInstanceMethod(pkg, chainCfg),
		generateRegisterMethod(elemName, pkg, chainCfg),
		generateChainElem(elemName, pkg, chainCfg),
	)
}

func generateInstanceMethod(pkg *types.Package, chainCfg chainConfig) Code {
	return Method(
		Receiver("c", Ptr(Id(chainCfg.ChainName))),
		"Instance",
		Sign(In(), Out(Param("", SmartQual(pkg.Name(), pkg.Path(), chainCfg.IfaceName), false))),
		Raw("return &c.root"),
	)
}

func generateRegisterMethod(elemName string, pkg *types.Package, chainCfg chainConfig) Code {
	return Method(
		Receiver("c", Ptr(Id(chainCfg.ChainName))),
		"Register",
		generateRegisterMethodSignature(pkg, chainCfg),
		generateRegisterMethodBody(elemName),
	)
}

func generateRegisterMethodSignature(pkg *types.Package, chainCfg chainConfig) Code {
	return Sign(
		In(
			Param("priority", Int, false),
			Param("f", FuncType(Sign(
				In(Param("", SmartQual(pkg.Name(), pkg.Path(), chainCfg.IfaceName), false)),
				Out(Param("", SmartQual(pkg.Name(), pkg.Path(), chainCfg.IfaceName), false)),
			)), false),
		),
		Out(),
	)
}

func generateRegisterMethodBody(elemName string) Code {
	return Lines(
		Raw("elem := &c.root"),
		For(Raw("elem.next != nil && elem.priority < priority"), Lines(
			Raw("elem = elem.next"),
		)),
		AssignAndDecl(Id("nextElem"), Inst(Addr(Id(elemName)), Fields(
			Field("priority", Raw("elem.priority")),
			Field("impl", Raw("elem.impl")),
			Field("next", Raw("elem.next")),
		))),
		Raw("elem.priority = priority"),
		Raw("elem.impl = f(nextElem)"),
		Raw("elem.next = nextElem"),
	)
}

func generateChainElem(elemName string, pkg *types.Package, chainCfg chainConfig) Code {
	return Blocks(
		Type(elemName, Struct(FieldDecls(
			FieldDecl("priority", Int),
			FieldDecl("impl", SmartQual(pkg.Name(), pkg.Path(), chainCfg.IfaceName)),
			FieldDecl("next", Ptr(Id(elemName))),
		))),
		generateChainElemMethods(elemName, chainCfg.Iface),
	)
}

func generateChainElemMethods(elemName string, iface *types.Interface) Code {
	methods := make([]Code, 0)
	tgen.ForEachInterfaceMethod(iface, func(name string, sign *types.Signature) {
		methods = append(methods, generateChainElemMethod(elemName, name, sign))
	})
	return Blocks(methods...)
}

func generateChainElemMethod(elemName string, name string, sign *types.Signature) Code {
	return Method(
		Receiver("e", Ptr(Id(elemName))),
		name,
		generateChainElemMethodSignature(sign),
		generateChainElemMethodBody(name, sign),
	)
}

func generateChainElemMethodSignature(sign *types.Signature) Code {
	params := sign.Params()
	paramsNames := utils.InIds(params.Len())
	n := params.Len()
	in := make([]Code, 0, n)
	if sign.Variadic() {
		n--
	}
	for i := 0; i < n; i++ {
		param := params.At(i)
		in = append(in, Param(paramsNames[i], tgen.Type(param.Type()), false))
	}
	if sign.Variadic() {
		param := params.At(n)
		in = append(in, Param(paramsNames[n], tgen.Type(param.Type()), true))
	}

	results := sign.Results()
	n = results.Len()
	out := make([]Code, 0, n)
	for i := 0; i < n; i++ {
		result := results.At(i)
		out = append(out, Param("", tgen.Type(result.Type()), false))
	}
	return Sign(In(in...), Out(out...))
}

func generateChainElemMethodBody(
	name string,
	sign *types.Signature,
) Code {
	params := sign.Params()
	paramsNames := utils.InIds(params.Len())
	results := sign.Results()
	implCall := Call(
		Join(Raw("e.impl."), Id(name)),
		Ids(paramsNames...),
	)
	if results.Len() == 0 {
		return implCall
	}
	return Return(implCall)
}
