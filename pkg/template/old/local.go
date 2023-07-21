package old

func readLocalBoxTemplate(cacheOpts *CacheSourceOpts, path string) (*BoxInfo, error) {

	if cacheOpts != nil {
		// TODO
		//return newCachedTemplateInfo(Local, "TODO_PATH")
	}

	if template, err := readBoxTemplate(path); err != nil {
		return nil, err
	} else {
		return &BoxTemplate{
			Template: template,
			Info:     newDefaultTemplateInfo(Local),
		}, nil
	}
}

func readLocalLabTemplate(cacheOpts *CacheSourceOpts, path string) (*LabInfo, error) {

	if cacheOpts != nil {
		// TODO
		//return newCachedTemplateInfo(Local, "TODO_PATH")
	}

	if template, err := readLabTemplate(path); err != nil {
		return nil, err
	} else {
		return &LabTemplate{
			Template: template,
			Info:     newDefaultTemplateInfo(Local),
		}, nil
	}
}
